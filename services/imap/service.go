package imap

import (
	"context"
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	_ "github.com/emersion/go-message/charset"
	"go.uber.org/zap"
)

type Service interface {
	Start(ctx context.Context, emailChan chan EmailEntity) error
	Fetch(emailNo uint32) ([]*EmailEntity, error)
}

func ProvideService(
	cfg *Cfg,
	log *zap.SugaredLogger,
) (srv Service, err error) {
	srv = &ServiceImpl{
		log: log,
		cfg: cfg,
	}
	return
}

type ServiceImpl struct {
	log        *zap.SugaredLogger
	cfg        *Cfg
	idleClient *imapclient.Client
}

func (srv *ServiceImpl) Start(ctx context.Context, emailChan chan EmailEntity) (err error) {
	srv.log.Infof("Starting Imap listener")
	idleC, err := imapclient.DialInsecure(fmt.Sprintf("%s:%d", srv.cfg.Host, srv.cfg.Port), &imapclient.Options{
		UnilateralDataHandler: &imapclient.UnilateralDataHandler{
			Expunge: func(seqNum uint32) {
				srv.log.Infof("message %v has been expunged", seqNum)
			},
			Mailbox: func(data *imapclient.UnilateralDataMailbox) {
				srv.log.Infof("mailbox: %+v", data)
				if data.NumMessages != nil {
					var emailEntities []*EmailEntity
					if emailEntities, err = srv.Fetch(*data.NumMessages); err != nil {
						srv.log.Warnf("Failed to fetch email: %v", err)
					}
					for _, emailEntity := range emailEntities {
						if emailEntity != nil {
							emailChan <- *emailEntity
						}
					}
				}
			},
		},
	})

	if err != nil {
		err = fmt.Errorf("failed to dial TLS: %v", err)
		return
	}

	srv.idleClient = idleC

	if err = idleC.Login(srv.cfg.EmailAddress, srv.cfg.Password).Wait(); err != nil {
		err = fmt.Errorf("failed to login: %v", err)
		return
	}

	if _, err = idleC.Select("INBOX", nil).Wait(); err != nil {
		err = fmt.Errorf("failed to select INBOX: %v", err)
		return
	}

	// Start idling
	for {
		srv.log.Infof("Starting idling")
		var idleCmd *imapclient.IdleCommand
		idleCmd, err = idleC.Idle()
		if err != nil {
			err = fmt.Errorf("failed to start idling: %v", err)
			return
		}
		if err = idleCmd.Wait(); err != nil {
			err = fmt.Errorf("failed to wait for idling: %v", err)
			return
		}
		if err = idleCmd.Close(); err != nil {
			return
		}
	}
}

// https://github.com/emersion/go-imap/issues/617
func (srv *ServiceImpl) Fetch(emailNo uint32) (emailEntities []*EmailEntity, err error) {
	srv.log.Infof("email no %+v", emailNo)
	// create client for fetching emails
	var fetchClient *imapclient.Client
	if fetchClient, err = imapclient.DialInsecure(fmt.Sprintf("%s:%d", srv.cfg.Host, srv.cfg.Port), nil); err != nil {
		err = fmt.Errorf("failed to dial insecure: %v", err)
		return
	}
	defer fetchClient.Close()
	if err = fetchClient.Login(srv.cfg.EmailAddress, srv.cfg.Password).Wait(); err != nil {
		err = fmt.Errorf("failed to login: %v", err)
		return
	}
	// select "INBOX" mailbox
	selectCmd := fetchClient.Select("INBOX", nil)
	var selectData *imap.SelectData
	if selectData, err = selectCmd.Wait(); err != nil {
		err = fmt.Errorf("failed to select INBOX: %v", err)
		return
	}
	if emailNo == 0 {
		emailNo = selectData.NumMessages
	}
	seqSet := imap.SeqSetNum(emailNo)
	fetchOptions := &imap.FetchOptions{
		UID:         true,
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{{}},
	}
	fetchCmd := fetchClient.Fetch(seqSet, fetchOptions)
	for {
		msg := fetchCmd.Next()
		if msg == nil {
			break
		}
		for {
			item := msg.Next()
			if item == nil {
				break
			}
			switch item := item.(type) {
			case imapclient.FetchItemDataUID:
				srv.log.Infof("UID: %v", item.UID)
			case imapclient.FetchItemDataBodySection:
				srv.log.Infof("section")
				var emailEntity *EmailEntity
				if emailEntity, err = ParseEmailLiteral(item.Literal); err != nil {
					srv.log.Errorf("Failed to parse email literal: %v", err)
				} else {
					srv.log.Info("add one email entity")
					emailEntities = append(emailEntities, emailEntity)
				}
			}
		}
	}
	return
}
