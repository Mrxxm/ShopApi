package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"shop_api/user-web/global"
)

// 1.从环境变量中获取配置
func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

// 2.初始化配置文件
func InitConfig() {
	debug := GetEnvInfo("SHOP")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug == "debug" {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}
	global.GetSugar().Infof("配置文件名称：%s", configFileName)

	vip := viper.New()
	// 设置配置文件路径
	vip.SetConfigFile(configFileName)
	if err := vip.ReadInConfig(); err != nil {
		panic(fmt.Errorf("无法读取配置文件: %s", err))
	}
	fmt.Println(vip.AllSettings(), vip.AllKeys())

	// 解析配置文件内容到结构体
	if err := vip.Unmarshal(&global.ServerConfig); err != nil {
		global.GetSugar().Errorf("无法解析配置文件内容: %v", err)
		panic(err)
	}

	// 打印解析后的内容
	global.GetSugar().Infof("配置文件内容：%+v", *global.ServerConfig)

	// 配置文件动态监控功能
	vip.WatchConfig() // 开始监控配置文件变化
	vip.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变化时触发的事件
		global.GetSugar().Infof("配置文件发生变化：%s", e.Name)
		// 重新读取配置文件并反序列化
		if err := vip.ReadInConfig(); err != nil {
			global.GetSugar().Errorf("重新读取配置文件失败: %v", err)
			return
		}
		if err := vip.Unmarshal(&global.ServerConfig); err != nil {
			global.GetSugar().Errorf("无法解析配置文件内容: %v", err)
			return
		}
		global.GetSugar().Infof("配置文件变化后的内容：%+v", *global.ServerConfig)
	})
}
