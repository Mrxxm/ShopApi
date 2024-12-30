package config

type ServerConfig struct {
	UserSrvConfig UserSrvConfig `mapstructure:"user_srv"`
	JWTConfig     JWTConfig     `mapstructure:"jwt"`
	ConsulConfig  ConsulConfig  `mapstructure:"consul"`
	Name          string        `mapstructure:"name"`
	Port          int           `mapstructure:"port"`
}

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
