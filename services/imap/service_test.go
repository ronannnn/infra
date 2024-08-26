package imap_test

import (
	"context"
	"testing"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/log"
	"github.com/ronannnn/infra/services/imap"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type ImapTestCfg struct {
	Imap cfg.Imap `mapstructure:"imap"`
	Log  cfg.Log  `mapstructure:"log"`
}

func TestImapFetchService(t *testing.T) {
	var err error
	testCfg := ImapTestCfg{}
	err = cfg.ReadFromFile("../../configs/config.imap.toml", &testCfg)
	require.NoError(t, err)
	// init logger
	var logger *zap.SugaredLogger
	logger, err = log.New(&testCfg.Log)
	require.NoError(t, err)
	// init imap service
	var srv imap.Service
	srv, err = imap.ProvideService(&testCfg.Imap, logger)
	require.NoError(t, err)
	_, err = srv.Fetch(0)
	require.NoError(t, err)
}

func TestImapStartService(t *testing.T) {
	var err error
	testCfg := ImapTestCfg{}
	err = cfg.ReadFromFile("../../configs/config.imap.toml", &testCfg)
	require.NoError(t, err)
	// init logger
	var logger *zap.SugaredLogger
	logger, err = log.New(&testCfg.Log)
	require.NoError(t, err)
	// init imap service
	var srv imap.Service
	srv, err = imap.ProvideService(&testCfg.Imap, logger)
	require.NoError(t, err)
	emailEntitiesChan := make(chan imap.EmailEntity)
	err = srv.Start(context.Background(), emailEntitiesChan)
	require.NoError(t, err)
	for {
		emailEntity := <-emailEntitiesChan
		logger.Infof("%+v", emailEntity)
	}
}
