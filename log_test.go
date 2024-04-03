package infra_test

import (
	"testing"
	"time"

	"github.com/ronannnn/infra"
	"github.com/ronannnn/infra/cfg"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLogRotatedFiles(t *testing.T) {
	var err error
	testCfg := cfg.Cfg{}
	err = cfg.ReadFromFile("configs/config.logtest.toml", &testCfg)
	require.NoError(t, err)
	// init log
	var log *zap.SugaredLogger
	log, err = infra.NewLog(&testCfg.Log)
	require.NoError(t, err)
	log.Info("test log")
	time.Sleep(10 * time.Second)
}
