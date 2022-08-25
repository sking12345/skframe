package config

import (
	"github.com/micro/go-micro"
	"skframe/pkg/config"
)

func init()  {
	config.Add("micro", func() map[string]interface{} {
		return map[string]interface{}{
			"addr":config.Env("MICRO_REGISTER_ADDR","127.0.0.1:8500"), //服务注册中心
			"name":config.Env("MICRO_SERVER_NAME","micro"),
			"version":config.Env("MICRO_SERVER_VERSION","v1.0"),
			"ttl":config.Env("MICRO_TTL",5), //单位s 注册服务的过期时间
			"interval":config.Env("MICRO_INTERVAL",5), //单位s,间隔多久再次注册服务
			"handler": func(service micro.Service) {  //注册服务

			},
		}
	})
}

