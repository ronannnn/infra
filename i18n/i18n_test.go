package i18n_test

import (
	"testing"

	"github.com/ronannnn/infra/i18n"
	"github.com/stretchr/testify/require"
)

type I18nWithData struct {
	Name  string
	Place string
}

func TestI18n(t *testing.T) {
	translator, err := i18n.New(i18n.Cfg{BundleDir: "./testdata/"})
	require.NoError(t, err)
	require.Equal(t, translator.Tr(i18n.LanguageChinese, "base.test"), "测试")
	require.Equal(t, translator.Tr(i18n.LanguageEnglish, "base.test"), "en_test")
	require.Equal(t, translator.TrWithData(i18n.LanguageChinese, "base.withData", I18nWithData{Name: "小陈", Place: "高新"}), "小陈在高新")
	require.Equal(t, translator.Tr(i18n.LanguageChinese, "base.not_found"), "base.not_found")
}
