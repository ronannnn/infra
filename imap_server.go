package infra

import (
	"fmt"

	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/ronannnn/infra/cfg"
)

func NewImapServer(cfg *cfg.Imap) (client *imapclient.Client, err error) {
	if client, err = imapclient.DialTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil); err != nil {
		return
	}
	if err = client.Login(cfg.EmailAddress, cfg.Password).Wait(); err != nil {
		return
	}
	return
}
