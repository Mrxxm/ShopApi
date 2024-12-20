package config

type UserSrvConfig struct {
	Host string `mapsstructure:"host"`
	Port int    `mapsstructure:"port"`
}

type ServerConfig struct {
	Name          string        `mapsstructure:"name"`
	Port          int           `mapsstructure:"port"`
	UserSrvConfig UserSrvConfig `mapsstructure:"user_srv"`
}
