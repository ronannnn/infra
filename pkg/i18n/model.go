package i18n

type LangType string

var DefaultLang = EnUs

const (
	ZhCn LangType = "zh-cn"
	EnUs LangType = "en-us"
)

type Cfg struct {
	LocalesDir             string `mapstructure:"locales-dir"`
	ZhCnTomlFilenamePrefix string `mapstructure:"zh-cn-toml-filename-prefix"`
	EnUsTomlFilenamePrefix string `mapstructure:"en-us-toml-filename-prefix"`
	CtxKey                 string `mapstructure:"ctx-key"` // ctxKey to get lang from context
}
