package cfg

type Imap struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	EmailAddress string `mapstructure:"email-address"`
	Password     string `mapstructure:"password"`
}
