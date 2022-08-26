package etcd
//
//import (
//	"context"
//	"errors"
//	"github.com/coreos/etcd/clientv3"
//	"github.com/coreos/etcd/mvcc/mvccpb"
//	"go.etcd.io/etcd/clientv3/concurrency"
//	"go.uber.org/zap"
//	"skframe/pkg/helpers"
//	"skframe/pkg/logger"
//	"sync"
//	"time"
//)
//
///*
// * 分布式锁
// */
//
//type SessionLock struct {
//	Session   *concurrency.Session
//	MutexLock *concurrency.Mutex
//}
//
//type EtcdClient struct {
//	Client      *clientv3.Client
//	KeyValList  map[string]string
//	Context     context.Context
//	SessionLock map[string]*SessionLock
//	WatchList   map[string]bool
//}
//
//var once sync.Once
//var Client *EtcdClient
//
//func ConnectEtcd(address string, dialTimeout int64) {
//	once.Do(func() {
//		Client = NewClient(address, dialTimeout)
//	})
//}
//
//func NewClient(address string, dialTimeout int64) *EtcdClient {
//	etcd := &EtcdClient{}
//	etcd.Context = context.Background()
//	client, err := clientv3.New(clientv3.Config{
//		Endpoints:   []string{address},
//		DialTimeout: time.Duration(dialTimeout) * time.Second,
//	})
//	logger.LogIf(err)
//	etcd.Client = client
//	etcd.KeyValList = map[string]string{}
//	etcd.SessionLock = map[string]*SessionLock{}
//	etcd.WatchList = map[string]bool{}
//	return etcd
//}
//
//func (etcd *EtcdClient) Close() {
//	for _, item := range etcd.SessionLock {
//		item.MutexLock.Unlock(etcd.Context)
//	}
//	for key, _ := range etcd.WatchList {
//		etcd.WatchList[key] = false
//	}
//	etcd.Client.Close()
//}
//
//func (etcd *EtcdClient) Exists(key string) bool { //新增,如果当前列表中不存在该数据，则向服务器请求一次数据
//	if helpers.Empty(etcd.KeyValList[key]) == true {
//		etcd.KeyValList[key] = etcd.Get(key)
//	}
//	return helpers.Empty(etcd.KeyValList[key])
//}
//
//func (etcd *EtcdClient) Set(key string, val string) bool { //新增
//	_, err := etcd.Client.Put(etcd.Context, key, val)
//	logger.LogIf(err)
//	if err == nil && helpers.Empty(etcd.WatchList[key]) == false && etcd.WatchList[key] == true {
//		etcd.KeyValList[key] = val
//	}
//	return err == nil
//}
//
//func (etcd *EtcdClient) Get(key string) (resultArr string) { //获取
//	if helpers.Empty(etcd.WatchList[key]) == false && etcd.WatchList[key] == true {
//		if helpers.Empty(etcd.KeyValList[key]) == false {
//			return etcd.KeyValList[key]
//		}
//	}
//	response, err := etcd.Client.Get(etcd.Context, key)
//	logger.LogIf(err)
//	for _, item := range response.Kvs {
//		resultArr = string(item.Value)
//	}
//	etcd.KeyValList[key] = resultArr
//	return
//}
//
//func (etcd *EtcdClient) Del(key string) bool { //删除
//	_, err := etcd.Client.Delete(etcd.Context, key)
//	logger.LogIf(err)
//	return err == nil
//}
//
//func (etcd *EtcdClient) Watch(key string, ChangeEvent func(key string, val string), DelEvent func(key string)) { //监听某个key
//	rChan := etcd.Client.Watch(etcd.Context, key)
//	etcd.WatchList[key] = true
//	for wResp := range rChan {
//		if etcd.WatchList[key] == false { //如果关闭该链接，则通关退出监听
//			break
//		}
//		for _, env := range wResp.Events {
//			switch env.Type {
//			case mvccpb.PUT:
//				ChangeEvent(string(env.Kv.Key), string(env.Kv.Value))
//				etcd.KeyValList[string(env.Kv.Key)] = string(env.Kv.Value)
//			case mvccpb.DELETE:
//				DelEvent(string(env.Kv.Key))
//				delete(etcd.KeyValList, string(env.Kv.Key))
//			}
//		}
//	}
//}
//
//func (etcd *EtcdClient) Lock(key string) bool { //开始锁
//	sessionId, err := concurrency.NewSession(etcd.Client)
//	if err != nil {
//		logger.Error("etcd", zap.Error(err))
//		return false
//	}
//	mutex := concurrency.NewMutex(sessionId, key)
//	err = mutex.Lock(etcd.Context)
//	if err != nil {
//		logger.Error("etcd", zap.Error(err))
//		return false
//	}
//	etcd.SessionLock[key] = &SessionLock{
//		MutexLock: mutex,
//		Session:   sessionId,
//	}
//	return true
//}
//
//func (etcd *EtcdClient) UnLock(key string) bool { //结束锁
//	if etcd.SessionLock[key] == nil {
//		logger.Error("etcd", zap.Error(errors.New(key+" not exists")))
//		return false
//	}
//	if err := etcd.SessionLock[key].MutexLock.Unlock(etcd.Context); err != nil {
//		logger.Error("etcd", zap.Error(err))
//		return false
//	}
//	delete(etcd.SessionLock, key)
//	return true
//}
