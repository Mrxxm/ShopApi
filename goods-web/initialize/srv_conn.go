package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/goods-web/global"
	"shop_api/goods-web/proto"
)

func InitSrvConn() {
	conn, err := grpc.Dial(fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port, global.ServerConfig.GoodsSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("连接 [user_srv] 失败", err.Error())
	}

	global.GoodsSrvClient = proto.NewGoodsClient(conn)
}
