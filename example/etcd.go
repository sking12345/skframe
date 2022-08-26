package example

func TestEtcdWatch() { //可多个监听
	//if etcd.Client.Exists("key") == true { //判断该值是否存在
	//	return
	//}
	//if etcd.Client.Set("key", "val") == false {
	//	return
	//}
	////if etcd.Etcd.Del("key") == false {
	////	return
	////}
	//arr := etcd.Client.Get("key")
	//fmt.Println(arr)
	//etcd.Client.Watch("key", func(s string, s2 string) { //change call
	//	fmt.Println(s, s2)
	//}, func(key string) { //delete call
	//	fmt.Println(key)
	//})
}

func TestLock() {
	//if etcd.Client.Lock("key1") == false {
	//	fmt.Println("lock fail")
	//	return
	//}
	//defer func() {
	//	if etcd.Client.UnLock("key1") == false {
	//		fmt.Println("unlock fail")
	//	}
	//	fmt.Println("unlock success")
	//}()
	//fmt.Println("lock success")
}
