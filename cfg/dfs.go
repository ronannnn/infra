package cfg

type Dfs struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key"`
	Location        string `mapstructure:"location"`
	RootBucket      string `mapstructure:"root-bucket"`
	Secure          bool   `mapstructure:"secure"`
}
