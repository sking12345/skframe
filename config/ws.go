package config

import "skframe/pkg/config"

func init() {
	config.Add("ws", func() map[string]interface{} {
		return map[string]interface{}{
			"path": config.Env("WS_PATH", "/"),
			"port": config.Env("WS_PORT", "3080"),
		}
	})
}
