package log

type Cfg struct {
	Level           string `mapstructure:"level"`
	StoreDir        string `mapstructure:"store-dir"`
	LatestFilename  string `mapstructure:"latest-filename"`
	TimeFormat      string `mapstructure:"time-format"`
	LogInConsole    bool   `mapstructure:"log-in-console"`
	LogInRotateFile bool   `mapstructure:"log-in-rotate-file"`
}
