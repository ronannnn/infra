package cfg

type Rabbitmq struct {
	Addr      string `mapstructure:"addr"`
	QueueName string `mapstructure:"queue-name"`
	EnableSsl bool   `mapstructure:"enable-ssl"`
}
