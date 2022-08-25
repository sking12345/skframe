package config

import "skframe/pkg/config"

func init() {
	config.Add("rabbit", func() map[string]interface{} {
		return map[string]interface{}{
			"url": config.Env("RABBIT_URL", ""),
		}
	})
}
