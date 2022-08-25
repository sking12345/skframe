package bootstrap

import (
	"skframe/pkg/config"
	"skframe/pkg/logger"
)

// SetupLogger 初始化 Logger
func SetupLogger() {

	logger.InitLogger(
		config.GetString("log.filename"),
		config.GetInt("log.max_size"),
		config.GetInt("log.max_backup"),
		config.GetInt("log.max_age"),
		config.GetBool("log.compress"),
		config.GetString("log.type"),
		config.GetString("log.level"),
	)
}
