package infra_test

import (
	"testing"

	"github.com/ronannnn/infra"
	"github.com/ronannnn/infra/cfg"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type LogTestCfg struct {
	Log cfg.Log `mapstructure:"log"`
}

func TestLogRotatedFiles(t *testing.T) {
	var err error
	testCfg := LogTestCfg{}
	err = cfg.ReadFromFile("configs/config.logtest.toml", &testCfg)
	require.NoError(t, err)
	// init log
	var log *zap.SugaredLogger
	log, err = infra.NewLog(&testCfg.Log)
	require.NoError(t, err)
	log.Info("test log")
}
