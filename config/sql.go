package config

import (
	"skframe/pkg/config"
)

func init() {

	config.Add("database", func() map[string]interface{} {
		return map[string]interface{}{
			// 默认数据库

			"master_slave": config.Env("DB_OPEN_MASTER_SLAVE", false),

			"mysql_master": map[string]interface{}{
				// 数据库连接信息
				"connection":      config.Env("DB_CONNECTION_MASTER", "mysql"),
				"host":            config.Env("DB_HOST_MASTER", "127.0.0.1"),
				"port":            config.Env("DB_PORT_MASTER", "3306"),
				"database":        config.Env("DB_DATABASE_MASTER", "test"),
				"username":        config.Env("DB_USERNAME_MASTER", ""),
				"password":        config.Env("DB_PASSWORD_MASTER", ""),
				"charset":         "utf8mb4",
				"max_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 5),
				// 连接池配置,gorm的配置
				//"max_idle_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 100),
				//"max_open_connections": config.Env("DB_MAX_OPEN_CONNECTIONS", 25),
				//"max_life_seconds":     config.Env("DB_MAX_LIFE_SECONDS", 5*60),
			},
			"mysql_save": map[string]interface{}{
				// 数据库连接信息
				"connection":      config.Env("DB_CONNECTION_SLAVE", "mysql"),
				"host":            config.Env("DB_HOST_SLAVE", "127.0.0.1"),
				"port":            config.Env("DB_PORT_SLAVE", "3306"),
				"database":        config.Env("DB_DATABASE_SLAVE", "test"),
				"username":        config.Env("DB_USERNAME_SLAVE", ""),
				"password":        config.Env("DB_PASSWORD_SLAVE", ""),
				"charset":         "utf8mb4",
				"max_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 5),
				// 连接池配置
				//"max_idle_connections": config.Env("DB_MAX_IDLE_CONNECTIONS", 100),
				//"max_open_connections": config.Env("DB_MAX_OPEN_CONNECTIONS", 25),
				//"max_life_seconds":     config.Env("DB_MAX_LIFE_SECONDS", 5*60),
			},
			"sqlite": map[string]interface{}{
				"database": config.Env("DB_SQL_FILE", "database/database.db"),
			},
		}
	})
}
