package main

import (
	"fmt"
	"shop_api/user-web/global"
	"shop_api/user-web/initialize"
)

func main() {
	var port int = 8022

	// 1.配置初始化日志
	initialize.Logger()
	global.GetSugar().Infof("启动服务器，端口：%d", port)

	// 2.初始化routers
	Router := initialize.Routers()

	// 3.启动服务
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		global.GetSugar().Panic("启动服务器失败:", err.Error())
	}
}
