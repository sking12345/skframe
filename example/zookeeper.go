package example

import (
	"fmt"
	"skframe/bootstrap"
	"skframe/pkg/zookeeper"
)

func TestZookeeper() {
	bootstrap.SetZookeeperUp()
	//zookeeper.Client.CreateEverNode("/test-node5", []byte("test2"))
	//list, _ := zookeeper.Client.GetChildrenList("/test-node")
	//fmt.Println(list)
	zookeeper.Client.Watch("/test-node5", func(event string, path string) bool {
		fmt.Println(event, path)
		return false
	})
}
