package infra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Jwt      JwtCfg `mapstructure:"jwt"`
	Username string `mapstructure:"username"`
}

func TestConfigReader(t *testing.T) {
	cfg := TestConfig{}
	err := ReadCfgFromFile("util_test.toml", &cfg)
	require.NoError(t, err)
	fmt.Printf("%+v\n", cfg)
	require.Equal(t, "ronan", cfg.Username)
	require.Equal(t, "abc123", cfg.Jwt.AccessTokenSecret)
}
