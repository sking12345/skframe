package config

import "skframe/pkg/config"

func init() {
	config.Add("es", func() map[string]interface{} {
		return map[string]interface{}{
			"url": config.Env("ES_URL", "http://127.0.0.1:9200/"),
		}
	})
}
