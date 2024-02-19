package infra

import (
	"testing"

	"github.com/ronannnn/infra/cfg"
	"github.com/stretchr/testify/require"
)

func TestI18n(t *testing.T) {
	cfg := cfg.I18n{
		LocalesDir:             "locales",
		ZhCnTomlFilenamePrefix: "zh-cn-test",
		EnUsTomlFilenamePrefix: "en-us-test",
		CtxKey:                 "lang",
	}
	// init service
	i18n := NewI18n(cfg)
	t.Run("Test zh-cn", func(t *testing.T) {
		require.Equal(t, "测试", i18n.T(LangTypeZhCn, "Test"))
	})
	t.Run("Test en-us", func(t *testing.T) {
		require.Equal(t, "test", i18n.T(LangTypeEnUs, "Test"))
	})
}
