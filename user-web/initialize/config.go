package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"shop_api/user-web/global"
)

// 1.从环境变量中获取配置
func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitConfig() {

	// 一、本地获取nacos配置
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
		panic(fmt.Errorf("无法读取配置文件: %s", err.Error()))
	}

	// 解析配置文件内容到结构体
	if err := vip.Unmarshal(&global.NacosConfig); err != nil {
		global.GetSugar().Errorf("无法解析配置文件内容: %v", err)
		panic(err)
	}

	// 打印解析后的内容
	global.GetSugar().Infof("配置文件内容：%+v", global.NacosConfig)

	// 二、nacos服务中获取其他配置
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Nacos.Host,
			Port:   uint64(global.NacosConfig.Nacos.Port),
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Nacos.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		//RotateTime:          "1h",
		//MaxAge:              3,
		LogLevel: "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		global.GetSugar().Info("创建nacos客户端异常：", err.Error())
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Nacos.DataId,
		Group:  global.NacosConfig.Nacos.Group})

	if err != nil {
		global.GetSugar().Info("获取nacos配置异常：", err.Error())
	}

	// json转换映射成结构体
	jsonBytesContent := []byte(content)
	err = json.Unmarshal(jsonBytesContent, &global.ServerConfig)
	if err != nil {
		global.GetSugar().Info("转换nacos配置异常：", err.Error())
	}
	global.GetSugar().Infof("转换nacos配置打印：%+v", global.ServerConfig)

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Nacos.DataId,
		Group:  global.NacosConfig.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
}

// 2.初始化配置文件
func InitConfig2() {
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
		panic(fmt.Errorf("无法读取配置文件: %s", err.Error()))
	}

	// 解析配置文件内容到结构体
	if err := vip.Unmarshal(&global.ServerConfig); err != nil {
		global.GetSugar().Errorf("无法解析配置文件内容: %v", err)
		panic(err)
	}

	// 打印解析后的内容
	global.GetSugar().Infof("配置文件内容：%+v", global.ServerConfig)

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
		global.GetSugar().Infof("配置文件变化后的内容：%+v", global.ServerConfig)
	})
}
