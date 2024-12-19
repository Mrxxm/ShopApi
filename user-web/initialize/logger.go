package initialize

import (
	"go.uber.org/zap"
	"shop_api/user-web/global"
)

func Logger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	global.Sugar = logger.Sugar()
}
