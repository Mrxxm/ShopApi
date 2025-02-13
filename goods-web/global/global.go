package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/goods-web/config"
	"shop_api/goods-web/proto"
)

var (
	ServerConfig   config.ServerConfig // 配置文件结构体
	Trans          ut.Translator       // 定义一个全局翻译器T
	GoodsSrvClient proto.GoodsClient   // 定义一个全局用户服务客户端

	NacosConfig config.NacosConfig
)
