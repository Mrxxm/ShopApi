package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"shop_api/goods-web/global"
	"shop_api/goods-web/initialize"
	"shop_api/goods-web/utils"
	"shop_api/goods-web/utils/register/consul"
)

func main() {
	// 1.初始化日志

	// 2.读取配置文件
	initialize.InitConfig()

	// 3.初始化routers
	Router := initialize.Routers()

	// 4.初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Panic("启动翻译器失败:", err.Error())
		return
	}
	// 5.初始化srv的链接
	initialize.InitSrvConn()

	// 本地开发环境端口号固定，线上环境启动获取端口号
	viper.AutomaticEnv()
	debug := viper.GetString("SHOP")
	if debug != "debug" {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// 6.注册验证器

	// 7.注册consul服务
	serviceID := fmt.Sprintf("%s", uuid.NewV4()) // 服务id
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)
	_ = register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceID)

	// 8.启动服务
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动服务器失败:", err.Error())
	}
}
