package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Registry struct {
	Host string
	Port int
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) DeRegister(serviceId string) error {

	// 1.初始化配置
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%v", r.Host, r.Port)

	// 2.创建一个consul客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	err = client.Agent().ServiceDeregister(serviceId)

	return err
}

func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	// 1.初始化配置
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%v", r.Host, r.Port)

	// 2.创建一个consul客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 3.注册服务&生成注册对象&生成检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           "http://192.168.15.21:8023/health",
		Timeout:                        "5s",  // 超时时间
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "10s", // 多久后注销服务
	}
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	return nil
}
