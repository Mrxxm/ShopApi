package config

type UserSrvConfig struct {
	Host string `mapsstructure:"host"`
	Port int    `mapsstructure:"port"`
}

type ServerConfig struct {
	UserSrvConfig UserSrvConfig `mapsstructure:"user_srv"`
	JWTConfig     JWTConfig     `mapsstructure:"jwt"`
	Name          string        `mapsstructure:"name"`
	Port          int           `mapsstructure:"port"`
}

type JWTConfig struct {
	SigningKey string `mapsstructure:"signing_key"`
}
