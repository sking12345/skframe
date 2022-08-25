package rpc

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"skframe/pkg/logger"
	"time"
)

type Micro struct {

}

func (*Micro)Start(addr,serverName,version string,ttl,interval uint,handler func(service micro.Service)) error {
	consulReg := consul.NewRegistry( //新建一个consul注册的地址，也就是我们consul服务启动的机器ip+端口
		registry.Addrs(addr),
	)
	server := micro.NewService( //创建一个新的服务对象
		micro.Name(serverName),
		micro.Version(version),
		micro.Registry(consulReg),
		micro.RegisterTTL(time.Duration(ttl)*time.Second),
		micro.RegisterInterval(time.Duration(interval)*time.Second),
	)
	server.Init()
	handler(server)
	err := server.Run()
	logger.LogIf(err)
	return err
}



