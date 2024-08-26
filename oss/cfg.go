package oss

type Cfg struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret"`
	ExpiredInSec    int64  `mapstructure:"expired-in-sec"`
	Location        string `mapstructure:"location"`
	RootBucket      string `mapstructure:"root-bucket"`
	Secure          bool   `mapstructure:"secure"`
}
