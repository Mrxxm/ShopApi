package initialize

import (
	"go.uber.org/zap"
	"shop_api/user-web/global"
)

func Logger() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	global.SetSugar(sugar)
}
