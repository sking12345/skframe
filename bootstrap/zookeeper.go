package bootstrap

import (
	"skframe/pkg/config"
	"skframe/pkg/zookeeper"
)

func SetZookeeperUp() {
	zookeeper.ConnectZookeeper(config.GetString("zk.addr"), config.GetInt64("zk.timeout"))
}
