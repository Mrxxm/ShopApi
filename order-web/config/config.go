package config

type ServerConfig struct {
	OrderSrvConfig     OrderSrvConfig     `mapstructure:"order_srv" json:"order_srv"`
	GoodsSrvConfig     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvConfig InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`

	AliPayInfo AliPayInfoConfig `mapstructure:"alipay" json:"alipay"`

	JWTConfig    JWTConfig    `mapstructure:"jwt" json:"jwt"`
	ConsulConfig ConsulConfig `mapstructure:"consul" json:"consul"`
	JaegerInfo   JaegerConfig `mapstructure:"consul" json:"jaeger"`
	Name         string       `mapstructure:"name" json:"name"`
	Port         int          `mapstructure:"port" json:"port"`
	Host         string       `mapstructure:"host" json:"host"`
	Tags         []string     `mapstructure:"tags" json:"tags"`
}

type AliPayInfoConfig struct {
	AppID        string `mapstructure:"appid" json:"appid"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
}

type OrderSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}
type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JaegerConfig struct {
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
