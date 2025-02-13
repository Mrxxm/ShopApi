package config

type ServerConfig struct {
	GoodsSrvConfig GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	JWTConfig      JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulConfig   ConsulConfig   `mapstructure:"consul" json:"consul"`
	Name           string         `mapstructure:"name" json:"name"`
	Port           int            `mapstructure:"port" json:"port"`
}

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key" json:"signing_key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Nacos Nacos `mapstructure:"nacos"`
}

type Nacos struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
