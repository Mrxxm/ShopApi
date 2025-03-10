package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/order-web/config"
	"shop_api/order-web/proto"
)

var (
	ServerConfig       config.ServerConfig   // 配置文件结构体
	Trans              ut.Translator         // 定义一个全局翻译器T
	OrderSrvClient     proto.OrderClient     // 定义一个全局用户服务客户端
	GoodsSrvClient     proto.GoodsClient     // 定义一个全局用户服务客户端
	InventorySrvClient proto.InventoryClient // 定义一个全局用户服务客户端

	NacosConfig config.NacosConfig
)
