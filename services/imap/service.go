package imap

import (
	"fmt"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	_ "github.com/emersion/go-message/charset"
	"github.com/ronannnn/infra/cfg"
	"go.uber.org/zap"
)

type Service interface {
	FetchLatestEmails() (*EmailEntity, error)
}

func ProvideService(
	cfg *cfg.Imap,
	log *zap.SugaredLogger,
) (srv Service, err error) {
	log.Info("Creating Imap service")
	defer func() {
		if err != nil {
			log.Error("Failed to create Imap service", err)
		} else {
			log.Info("Imap service created")
		}
	}()
	var client *imapclient.Client
	if client, err = imapclient.DialInsecure(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil); err != nil {
		return
	}
	if err = client.Login(cfg.EmailAddress, cfg.Password).Wait(); err != nil {
		return
	}
	srv = &ServiceImpl{
		client: client,
		log:    log,
	}
	return
}

type ServiceImpl struct {
	client *imapclient.Client
	log    *zap.SugaredLogger
}

// https://github.com/emersion/go-imap/issues/617
func (srv *ServiceImpl) FetchLatestEmails() (emailEntity *EmailEntity, err error) {
	selectCmd := srv.client.Select("INBOX", nil)
	var selectedData *imap.SelectData
	if selectedData, err = selectCmd.Wait(); err != nil {
		return
	}
	seqSet := imap.SeqSetNum(selectedData.NumMessages)
	fetchOptions := &imap.FetchOptions{
		UID:         true,
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{{}},
	}
	fetchCmd := srv.client.Fetch(seqSet, fetchOptions)
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
				if emailEntity, err = ParseEmailLiteral(item.Literal); err != nil {
					srv.log.Errorf("Failed to parse email literal: %v", err)
					continue
				}
				srv.log.Infof("EmailEntity: %v", emailEntity.String())
			}
		}
	}
	return
}
