package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shop_api/order-web/global"
	"shop_api/order-web/initialize"
	"shop_api/order-web/utils"
	"shop_api/order-web/utils/register/consul"
	"syscall"
)

func main() {
	// 1.初始化日志
	logger, _ := zap.NewDevelopment() // 默认包含Debug级别
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

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
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动服务器失败:", err.Error())
		} // 阻塞的方法,需要放在goroutine中，否则后续代码无法执行
	}()

	// 9.优雅退出，接收终止信号
	quit := make(chan os.Signal) // 无缓冲区通道
	// SIGINT（通常是用户按下 Ctrl+C）和 SIGTERM（通常是终止进程的信号）。当程序收到这两个信号之一时，操作系统会将信号发送到 quit 通道。
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 它用于接收来自操作系统的信号（比如中断信号 SIGINT 或终止信号 SIGTERM）
	<-quit                                               // 这一行代码会阻塞程序的执行，直到从 quit 通道接收到信号
	if err := register_client.DeRegister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
