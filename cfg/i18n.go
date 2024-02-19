package cfg

type I18n struct {
	LocalesDir             string `mapstructure:"locales-dir"`
	ZhCnTomlFilenamePrefix string `mapstructure:"zh-cn-toml-filename-prefix"`
	EnUsTomlFilenamePrefix string `mapstructure:"en-us-toml-filename-prefix"`
	CtxKey                 string `mapstructure:"ctx-key"` // ctxKey to get lang from context
}
