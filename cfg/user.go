package cfg

type User struct {
	DefaultHashedPassword string `mapstructure:"default-hashed-password"`
}
