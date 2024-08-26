package log_test

import (
	"testing"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type LogTestCfg struct {
	Log log.Cfg `mapstructure:"log"`
}

func TestLogRotatedFiles(t *testing.T) {
	var err error
	testCfg := LogTestCfg{}
	err = cfg.ReadFromFile("configs/config.logtest.toml", &testCfg)
	require.NoError(t, err)
	// init log
	var logger *zap.SugaredLogger
	logger, err = log.New(&testCfg.Log)
	require.NoError(t, err)
	logger.Info("test log")
}
