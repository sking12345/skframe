package rpc

import (
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"skframe/pkg/config"
	"skframe/pkg/logger"
	"time"
)

type Micro struct {

}



func (*Micro)NewServer(addr,name string,handler func(service micro.Service)) bool  {
	if handler == nil {
		logger.Error("micro server:",zap.Any("faild","handler is nil"))
		return false
	}
	consulReg := consul.NewRegistry(registry.Addrs(addr))
	server := micro.NewService(
		micro.Name(name),
		micro.Registry(consulReg),
		micro.RegisterTTL(time.Duration(config.GetInt("micro.ttl",5))*time.Second),
		micro.RegisterInterval(time.Duration(config.GetInt("micro.interval",10))*time.Second),
		)
	server.Init()
	handler(server)
	err := server.Run()
	if err != nil {
		logger.Error("micro server",zap.Any("run",err))
		return  false
	}
	return  true
}

func (*Micro)NewClient(addr string,handler func(service micro.Service))  bool {
	if handler == nil {
		logger.Error("micro server:",zap.Any("faild","handler is nil"))
		return false
	}
	consulReg := consul.NewRegistry(registry.Addrs(addr))
	server := micro.NewService(
		micro.Registry(consulReg),
		micro.RegisterTTL(time.Duration(config.GetInt("micro.ttl",5))*time.Second),
		micro.RegisterInterval(time.Duration(config.GetInt("micro.interval",10))*time.Second),
	)
	server.Init()
	handler(server)
	return true
}

