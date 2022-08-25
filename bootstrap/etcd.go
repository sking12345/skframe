package bootstrap

import (
	"skframe/pkg/config"
	"skframe/pkg/etcd"
)

func SetUpEtcd() {
	etcd.ConnectEtcd(
		config.GetString("etcd.host"),
		config.GetInt64("etcd.dial_timeout"),
	)
}
