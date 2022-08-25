package config

import (
	"skframe/pkg/config"
)

func init() {
	config.Add("etcd", func() map[string]interface{} {
		return map[string]interface{}{
			"host":         config.Env("ETCD_HOST", "127.0.0.1:2379"),
			"dial_timeout": config.Env("ETCD_DIAL_TIMEOUT", 1), //分钟
		}
	})
}
