package email

type Cfg struct {
	SmtpAddr     string `mapstructure:"smtp-addr"`
	SmtpPort     uint16 `mapstructure:"smtp-port"`
	EmailAccount string `mapstructure:"email-account"`
	EmailPasswd  string `mapstructure:"email-passwd"`
}
