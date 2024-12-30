package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"shop_api/user-web/global"
	"shop_api/user-web/initialize"
	"shop_api/user-web/utils"
	myvalidator "shop_api/user-web/validator"
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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// 7.启动服务
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		global.GetSugar().Panic("启动服务器失败:", err.Error())
	}
}
