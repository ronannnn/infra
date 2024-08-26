package cfg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Username string `mapstructure:"username"`
}

func TestConfigReader(t *testing.T) {
	cfg := TestConfig{}
	err := ReadFromFile("util_test.toml", &cfg)
	require.NoError(t, err)
	fmt.Printf("%+v\n", cfg)
	require.Equal(t, "ronan", cfg.Username)
}
