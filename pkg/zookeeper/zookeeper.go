package zookeeper

import (
	"errors"
	zk "github.com/samuel/go-zookeeper/zk"
	"go.uber.org/zap"
	"skframe/pkg/console"
	"skframe/pkg/logger"
	"sync"
	"time"
)

type Zookeeper struct {
	conn *zk.Conn
}

var Client *Zookeeper
var once sync.Once

func ConnectZookeeper(addr string, dialTimeout int64) {
	once.Do(func() {
		Client = NewClient(addr, dialTimeout, nil)
	})
}

func NewClient(addr string, dialTimeout int64, watchHandler func(event zk.Event)) *Zookeeper {
	var err error
	var con *zk.Conn
	if watchHandler == nil {
		con, _, err = zk.Connect([]string{addr}, time.Duration(dialTimeout)*time.Second, zk.WithEventCallback(watchHandler))
	} else {
		con, _, err = zk.Connect([]string{addr}, time.Duration(dialTimeout)*time.Second)
	}
	if err != nil {
		logger.Error("zk", zap.Error(err))
		return nil
	}
	console.Info("zk connect success")
	return &Zookeeper{
		conn: con,
	}
}

func (ptr *Zookeeper) Close() {
	ptr.conn.Close()
}

func (ptr *Zookeeper) CreateEverNode(path string, data []byte) bool { //创建永久节点，手动删除
	_, err := ptr.conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		logger.Error("zk", zap.Error(err))
	}
	return err == nil
}

func (ptr *Zookeeper) CreateTempNode(path string, data []byte) bool { //创建临时,连接断开后，自动删除
	_, err := ptr.conn.Create(path, data, 1, zk.WorldACL(zk.PermAll))
	if err != nil {
		logger.Error("zk", zap.Error(err))
	}
	return err == nil
}

func (ptr *Zookeeper) GetData(path string) (val string, version int32) {
	res, state, _ := ptr.conn.Get(path)
	return string(res), state.Version
}

func (ptr *Zookeeper) Exists(path string) (exits bool) {
	state, _, _ := ptr.conn.Exists(path)
	return state
}

func (ptr *Zookeeper) Del(path string, version int32) bool {
	err := ptr.conn.Delete(path, version)
	if err != nil {
		logger.Error("zk", zap.Error(err))
		return false
	}
	return true
}

func (ptr *Zookeeper) Watch(path string, handler func(event string, path string) bool) bool {
	if handler == nil {
		logger.Error("zk", zap.Error(errors.New("zk watch callFunc handler is null")))
		return false
	}
	for {
		state, _, event, err := ptr.conn.ExistsW(path)
		if err != nil {
			logger.Error("zk", zap.Error(err))
			return false
		}
		if state == false {
			return false
		}
		evt := <-event
		if handler(evt.Type.String(), evt.Path) == false {
			return true
		}
	}

}

func (ptr *Zookeeper) Update(path string, newData []byte, version int32) bool {
	_, err := ptr.conn.Set(path, newData, version)
	if err != nil {
		logger.Error("zk", zap.Error(err))
		return false
	}
	return true
}

func (ptr *Zookeeper) GetChildrenList(path string) (list []string, version int32) {
	list, state, _ := ptr.conn.Children(path)
	return list, state.Version
}
