package imap_test

import (
	"testing"

	"github.com/ronannnn/infra"
	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/services/imap"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type ImapTestCfg struct {
	Imap cfg.Imap `mapstructure:"imap"`
	Log  cfg.Log  `mapstructure:"log"`
}

func TestImapService(t *testing.T) {
	var err error
	testCfg := ImapTestCfg{}
	err = cfg.ReadFromFile("../../configs/config.imap.toml", &testCfg)
	require.NoError(t, err)
	// init log
	var log *zap.SugaredLogger
	log, err = infra.NewLog(&testCfg.Log)
	require.NoError(t, err)
	// init imap service
	var srv imap.Service
	srv, err = imap.ProvideService(&testCfg.Imap, log)
	require.NoError(t, err)
	_, err = srv.FetchLatestEmails()
	require.NoError(t, err)
}
