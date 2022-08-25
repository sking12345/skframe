package config

import "skframe/pkg/config"

func init() {
	config.Add("zk", func() map[string]interface{} {
		return map[string]interface{}{
			"addr":    config.Env("ZK_ADDR", "127.0.0.1:2181"),
			"timeout": config.Env("ZK_TIMEOUT", 5),
		}
	})

}
