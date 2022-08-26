package config

import (
	"skframe/example"
	"skframe/pkg/config"
)


func init()  {
	config.Add("micro", func() map[string]interface{} {
		return map[string]interface{}{
			"ttl":config.Env("MICRO_TTL",5), //单位s 注册服务的过期时间
			"interval":config.Env("MICRO_INTERVAL",5), //单位s,间隔多久再次注册服务
			"handler": func() {
				example.MicroClient()
			},
		}
	})
}

