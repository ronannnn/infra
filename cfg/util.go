package cfg

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func ReadFromFile(configFilepath string, cfg any) (err error) {
	v := viper.New()
	v.SetConfigFile(configFilepath)
	v.SetConfigType("toml")
	if err = v.ReadInConfig(); err != nil {
		return
	}
	if err = v.Unmarshal(cfg); err != nil {
		return
	}
	v.WatchConfig()
	// watching and updating Conf without application restart
	// TODO restart resources such as db
	v.OnConfigChange(func(e fsnotify.Event) {
		if err = v.Unmarshal(cfg); err != nil {
			panic(err)
		}
	})
	return
}
