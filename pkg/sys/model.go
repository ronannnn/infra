package sys

type Cfg struct {
	HttpAddr string `mapstructure:"http-addr"`
	HttpPort int    `mapstructure:"http-port"`
}
