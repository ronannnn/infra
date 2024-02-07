package i18n

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestI18n(t *testing.T) {
	cfg := Cfg{
		LocalesDir:             "locales",
		ZhCnTomlFilenamePrefix: "zh-cn-test",
		EnUsTomlFilenamePrefix: "en-us-test",
		CtxKey:                 "lang",
	}
	// init service
	i18n := New(cfg)
	t.Run("Test zh-cn", func(t *testing.T) {
		require.Equal(t, "测试", i18n.T(ZhCn, "Test"))
	})
	t.Run("Test en-us", func(t *testing.T) {
		require.Equal(t, "test", i18n.T(EnUs, "Test"))
	})
}
