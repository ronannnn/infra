package infra

type SysCfg struct {
	HttpAddr string `mapstructure:"http-addr"`
	HttpPort int    `mapstructure:"http-port"`
}
