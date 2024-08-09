package infra

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ronannnn/infra/cfg"
	"go.uber.org/zap"
)

// example from https://github.com/rabbitmq/amqp091-go/blob/main/_examples/client/client.go
// RabbitmqClient is the base struct for handling connection recovery, consumption and
// publishing. Note that this struct has an internal mutex to safeguard against
// data races. As you develop and iterate over this example, you may need to add
// further locks, or safeguards, to keep your application safe from data races
type RabbitmqClient struct {
	log             *zap.SugaredLogger
	m               *sync.Mutex
	queueName       string
	connection      *amqp.Connection
	channel         *amqp.Channel
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	done            chan bool
	isReady         bool
}

const (
	// When reconnecting to the server after connection failure
	reconnectDelay = 5 * time.Second

	// When setting up the channel after a channel exception
	reInitDelay = 2 * time.Second

	// When resending messages the server didn't confirm
	resendDelay = 5 * time.Second
)

var (
	errNotConnected  = fmt.Errorf("not connected to a server")
	errAlreadyClosed = fmt.Errorf("already closed: not connected to the server")
	errShutdown      = fmt.Errorf("client is shutting down")
)

// New creates a new consumer state instance, and automatically
// attempts to connect to the server.
func NewRabbitMq(
	log *zap.SugaredLogger,
	rmqCfg *cfg.Rabbitmq,
) *RabbitmqClient {
	client := RabbitmqClient{
		m:         &sync.Mutex{},
		log:       log,
		queueName: rmqCfg.QueueName,
		done:      make(chan bool),
	}
	go client.handleReconnect(rmqCfg.Addr)
	return &client
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (client *RabbitmqClient) handleReconnect(addr string) {
	for {
		client.m.Lock()
		client.isReady = false
		client.m.Unlock()

		client.log.Info("[rmq] attempting to connect")

		conn, err := client.connect(addr)
		if err != nil {
			client.log.Error("[rmq] failed to connect. Retrying...")

			select {
			case <-client.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := client.handleReInit(conn); done {
			break
		}
	}
}

// connect will create a new AMQP connection
func (client *RabbitmqClient) connect(addr string) (conn *amqp.Connection, err error) {
	if conn, err = amqp.Dial(addr); err != nil {
		return nil, err
	}

	client.changeConnection(conn)
	client.log.Info("[rmq] connected")
	return
}

// handleReInit will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (client *RabbitmqClient) handleReInit(conn *amqp.Connection) bool {
	for {
		client.m.Lock()
		client.isReady = false
		client.m.Unlock()

		err := client.init(conn)
		if err != nil {
			client.log.Error("[rmq] failed to initialize channel, retrying...")

			select {
			case <-client.done:
				return true
			case <-client.notifyConnClose:
				client.log.Info("[rmq] connection closed, reconnecting...")
				return false
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-client.done:
			return true
		case <-client.notifyConnClose:
			client.log.Info("[rmq] connection closed, reconnecting...")
			return false
		case <-client.notifyChanClose:
			client.log.Info("[rmq] channel closed, re-running init...")
		}
	}
}

// init will initialize channel & declare queue
func (client *RabbitmqClient) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		client.queueName,
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	client.changeChannel(ch)
	client.m.Lock()
	client.isReady = true
	client.m.Unlock()
	client.log.Info("[rmq] client init done")

	return nil
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (client *RabbitmqClient) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (client *RabbitmqClient) changeChannel(channel *amqp.Channel) {
	client.channel = channel
	client.notifyChanClose = make(chan *amqp.Error, 1)
	client.notifyConfirm = make(chan amqp.Confirmation, 1)
	client.channel.NotifyClose(client.notifyChanClose)
	client.channel.NotifyPublish(client.notifyConfirm)
}

// Push will push data onto the queue, and wait for a confirmation.
// This will block until the server sends a confirmation. Errors are
// only returned if the push action itself fails, see UnsafePush.
func (client *RabbitmqClient) Push(data []byte) error {
	for {
		client.m.Lock()
		if client.isReady {
			client.m.Unlock()
			break
		}
		client.m.Unlock()
		client.log.Warn("[rmq] not connected, waiting...")
		time.Sleep(reconnectDelay)
	}
	for {
		err := client.UnsafePush(data)
		if err != nil {
			client.log.Error("[rmq] push failed. Retrying...")
			select {
			case <-client.done:
				return errShutdown
			case <-time.After(resendDelay):
			}
			continue
		}
		confirm := <-client.notifyConfirm
		if confirm.Ack {
			client.log.Infof("[rmq] push confirmed [%d]", confirm.DeliveryTag)
			return nil
		}
	}
}

// UnsafePush will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (client *RabbitmqClient) UnsafePush(data []byte) error {
	client.m.Lock()
	if !client.isReady {
		client.m.Unlock()
		return errNotConnected
	}
	client.m.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return client.channel.PublishWithContext(
		ctx,
		"",               // Exchange
		client.queueName, // Routing key
		false,            // Mandatory
		false,            // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

// Consume will continuously put queue items on the channel.
// It is required to call delivery.Ack when it has been
// successfully processed, or delivery.Nack when it fails.
// Ignoring this will cause data to build up on the server.
func (client *RabbitmqClient) Consume() (<-chan amqp.Delivery, error) {
	for {
		client.m.Lock()
		if client.isReady {
			client.m.Unlock()
			break
		}
		client.m.Unlock()
		client.log.Warn("[rmq] not connected, waiting...")
		time.Sleep(reconnectDelay)
	}

	if err := client.channel.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	); err != nil {
		return nil, err
	}

	return client.channel.Consume(
		client.queueName,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
}

// Close will cleanly shut down the channel and connection.
func (client *RabbitmqClient) Close() error {
	client.m.Lock()
	// we read and write isReady in two locations, so we grab the lock and hold onto
	// it until we are finished
	defer client.m.Unlock()

	if !client.isReady {
		return errAlreadyClosed
	}
	close(client.done)
	err := client.channel.Close()
	if err != nil {
		return err
	}
	err = client.connection.Close()
	if err != nil {
		return err
	}

	client.isReady = false
	return nil
}
