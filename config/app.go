package config

import (
	"skframe/pkg/config"
)

func init()  {
	config.Add("app", func() map[string]interface{} {

		return map[string]interface{}{
			"name": config.Env("APP_NAME", "skFrame"),
			// 当前环境，用以区分多环境，一般为 local, stage, production, test
			"env":   config.Env("APP_ENV", "production"),
			"debug": config.Env("APP_DEBUG", false),
		}
	})
}