package infra

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { // path exists
		return true, nil
	} else if os.IsNotExist(err) { // error is 'not exist'
		return false, nil
	}
	return false, err // other error
}

func createDirsIfNotExist(dirs ...string) (err error) {
	for _, dir := range dirs {
		if existing, pathExistsErr := pathExists(dir); !existing && pathExistsErr == nil {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return
			}
		}
	}
	return
}

func ReadCfgFromFile(configFilepath string, cfg any) (err error) {
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
