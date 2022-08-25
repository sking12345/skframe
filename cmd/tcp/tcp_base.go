package tcp

type MsgProtocol struct {
	size int
	data []byte
}

type Server struct {
	SocketFd            int
	Port                int
	newConnectHandler   func(fd int)
	newMessageHandler   func(fd int, data []byte)
	closeConnectHandler func(fd int)
}

func NewTcp(port int) *Server {
	return &Server{Port: port}
}

func (ptr *Server) SetNewConnectHandler(handler func(fd int)) *Server {
	ptr.newConnectHandler = handler
	return ptr
}
func (ptr *Server) SetNewMessageHandler(handler func(fd int, data []byte)) *Server {
	ptr.newMessageHandler = handler
	return ptr
}
func (ptr *Server) SetCloseConnectHandler(handler func(fd int)) *Server {
	ptr.closeConnectHandler = handler
	return ptr
}
