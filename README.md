#skframe 是一个基与golang 设计支持try.catch 开发的业务框架，重构了linux/mac 的tcp,dup,websocket 的网络通信,避免一个连接就形成一个常驻的协程所造成不必要的资源消耗
#支持同一套业务对 tcp,udp,http,rpc 的全面支持 