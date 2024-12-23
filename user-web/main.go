package main

import (
	"fmt"
	"shop_api/user-web/global"
	"shop_api/user-web/initialize"
)

func main() {
	// 1.初始化日志
	initialize.Logger()

	// 2.读取配置文件
	initialize.InitConfig()

	// 3.初始化routers
	Router := initialize.Routers()

	// 4.初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		global.GetSugar().Panic("启动翻译器失败:", err.Error())
		return
	}

	// 5.启动服务
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		global.GetSugar().Panic("启动服务器失败:", err.Error())
	}
}
