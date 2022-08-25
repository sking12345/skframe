package bootstrap

import (
	"fmt"
	"skframe/pkg/cache"
	"skframe/pkg/config"
)

// SetupRedis 初始化 Redis

func SetupCache() {
	// 建立 Redis 连接
	cache.ConnectRedis(
		fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		config.GetString("redis.username"),
		config.GetString("redis.password"),
		config.GetInt("redis.database"),
	)
}
