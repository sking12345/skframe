package udp

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"skframe/pkg/logger"
	"skframe/pkg/try"
	"syscall"
	"unsafe"
)

/*
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

int recvfromUdp(int fd,char *dataBuf, int dataSize,char *clientAddr){
	socklen_t len = 16;
	return recvfrom(fd,dataBuf,dataSize,0,(struct sockaddr*)clientAddr,&len);
}
size_t sendData(int fd,char *data,int dataSize, char *addr){
	socklen_t len = 16;
	size_t wsize = sendto(fd,data,dataSize,0,(struct sockaddr *)addr,len);
	return wsize;
}
*/
import "C"

type Server struct {
	port          int
	buffMaxSize   int
	newMsgHandler func(fd int, data []byte, addr []byte)
}

func NewServer(port, buffMaxSize int) *Server {
	return &Server{
		port:        port,
		buffMaxSize: buffMaxSize,
	}
}
func SendData(fd int, data []byte, addr []byte) bool {
	status := C.sendData(C.int(fd), (*C.char)(unsafe.Pointer(&data[0])), C.int(len(data)), (*C.char)(unsafe.Pointer(&addr[0])))
	if status <= 0 {
		return false
	}
	return true
}

func (ptr *Server) SetMessageHandler(handler func(fd int, data []byte, addr []byte)) *Server {
	ptr.newMsgHandler = handler
	return ptr
}
func (ptr *Server) Run() {
	try.NewTry(func() {
		fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
		if err != nil {
			panic(err)
		}
		defer syscall.Close(fd)
		addr := syscall.SockaddrInet4{Port: ptr.port}
		copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
		err = syscall.Bind(fd, &addr)
		if err != nil {
			panic(err)
		}
		logger.Info("udp", zap.Any("status", "success"), zap.Any("port", ptr.port))
		for {
			dataBuf := make([]byte, ptr.buffMaxSize)
			ClientAddr := make([]byte, 16)
			size := C.recvfromUdp(C.int(fd), (*C.char)(unsafe.Pointer(&dataBuf[0])), C.int(ptr.buffMaxSize), (*C.char)(unsafe.Pointer(&ClientAddr[0])))
			if ptr.newMsgHandler != nil {
				go ptr.newMsgHandler(fd, dataBuf[:size], ClientAddr)
			}
		}
	}).Catch(func(err interface{}) {
		logger.ErrorString("udp", "run err:", fmt.Sprintf("%s", err))
		os.Exit(1)
	}).Run()
}
