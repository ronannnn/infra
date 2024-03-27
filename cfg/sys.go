package cfg

type Sys struct {
	HttpAddr    string `mapstructure:"http-addr"`
	HttpPort    int    `mapstructure:"http-port"`
	PrintRoutes bool   `mapstructure:"print-routes"`
}
