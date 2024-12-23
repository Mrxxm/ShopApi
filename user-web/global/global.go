package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/user-web/config"
)

var (
	ServerConfig *config.ServerConfig // 配置文件结构体
	Trans        ut.Translator        // 定义一个全局翻译器T
)
