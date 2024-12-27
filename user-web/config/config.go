package config

type ServerConfig struct {
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv"`
	JWTConfig     JWTConfig     `mapstructure:"jwt"`
	Name          string        `mapstructure:"name"`
	Port          int           `mapstructure:"port"`
}

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"`
}
