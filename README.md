## skframe 是一个基与golang 设计支持try.catch 开发的性能稳定而快速开发的业务框架，
## 重构了linux/mac 的tcp,dup 的网络通信,避免一个连接就形成一个常驻的协程所造成不必要的资源消耗
## 支持同一套业务对 tcp,udp,http,rpc 的全面支持 
## 支持命令创建model (go run main.go make model table_name)

## 启动服务
### 1: go run main.go http (默认)
### 2: go run main.go tcp 
### 3: go run main.go udp
### 4: go run main.go websocket
### 5: go run main.go rpc

## 待处理:
### 1:windows环境下 对 tcp,udp 支持的重构， 
### 2:websocket 的重构
### 3: 详细文档制作

