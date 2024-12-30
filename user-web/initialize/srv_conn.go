package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"shop_api/user-web/global"
	"shop_api/user-web/proto"
)

func InitSrvConn() {
	// 服务发现-1.初始化配置
	userSrvHost := ""
	userSrvPort := 0
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)

	// 服务发现-2.创建一个consul客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 服务发现-3.获取所有服务
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvConfig.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	if userSrvHost == "" {
		global.GetSugar().Fatal("用户服务不可达")
		return
	}

	// 4.拨号连接grpc服务
	connect, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		global.GetSugar().Errorw("[GetUserList] 连接 [user_srv] 失败", "msg", err.Error())
	}
	//defer connect.Close()
	// 问题：1.后续用户服务下线 2.改端口 3.改ip - 负载均衡来处理
	// 5.初始化全局变量(已经实现创立好了连接，这样后续就不用再进行TCP三次握手)
	// 问题：一个连接多个goroutine共用，性能会有影响 - 连接池
	userSrvClient := proto.NewUserClient(connect)
	global.UserSrvClient = userSrvClient
}
