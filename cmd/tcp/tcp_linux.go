//+build linux

package tcp

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
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

struct MsgProtocol {
    int size;
    char buf[0];
};

void destruct1(char **buf){
	if (*buf != NULL){
		free(*buf);
		*buf= NULL;
	}
}
char *newMsg(int fd){
	int size = 0;
	size_t readSizze = read(fd,&size,sizeof(int));
	printf("%ld\n",readSizze);
	if (readSizze <= 0){
		return NULL;
	}
	char *buf = NULL;
	buf = (char*)malloc(size);
	read(fd,buf,size);
	return buf;
}

size_t send(int fd,const char *data,int dataSize){
	struct MsgProtocol *msg = NULL;
	msg = (struct MsgProtocol*)malloc(sizeof(struct MsgProtocol)+dataSize);
	msg->size = dataSize;
	memcpy(msg->buf,data,dataSize);
	size_t wsize = write(fd,msg,sizeof(struct MsgProtocol)+dataSize);
	destruct1((char**)&msg);
	return wsize;

}

void TestPtr(char *buf){
printf("%p\n",buf);
}

*/
import "C"

const (
	EPOLLET        = 1 << 31
	MaxEpollEvents = 32
)

func SendMessage(fd int, data []byte) bool {
	sendStatus := C.send(C.int(fd), (*C.char)(unsafe.Pointer(&data)), C.int(len(data)))
	if int(sendStatus) != len(data)+4 {
		logger.Info("tcp", zap.Any("send error", sendStatus))
		return false
	}
	return true
}

func (ptr *Server) Run() {
	try.NewTry(func() {
		var event syscall.EpollEvent
		var events [MaxEpollEvents]syscall.EpollEvent
		fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
		if err != nil {
			panic(err)
		}
		defer syscall.Close(fd)
		if err = syscall.SetNonblock(fd, true); err != nil {
			panic(err)
		}
		addr := syscall.SockaddrInet4{Port: ptr.Port}
		copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
		syscall.Bind(fd, &addr)
		syscall.Listen(fd, MaxEpollEvents) //backlog:同时处理最大连接数
		epfd, e := syscall.EpollCreate1(0)
		if e != nil {
			panic(e)
		}
		defer syscall.Close(epfd)
		event.Events = syscall.EPOLLIN
		event.Fd = int32(fd)
		if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
			panic(e)
		}
		logger.Info("tcp", zap.String("status", "success"), zap.Any("port", ptr.Port))
		for {
			nevents, e := syscall.EpollWait(epfd, events[:], -1)
			if e != nil {
				logger.Info("tcp->epollWait", zap.Error(err))
				break
			}
			for ev := 0; ev < nevents; ev++ {
				if int(events[ev].Fd) == fd {
					connFd, _, err := syscall.Accept(fd)
					if err != nil {
						logger.Info("tcp->accept", zap.Error(err))
						continue
					}
					syscall.SetNonblock(connFd, true)
					event.Events = syscall.EPOLLIN | EPOLLET
					event.Fd = int32(connFd)
					if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connFd, &event); err != nil {
						logger.Info("tcp->epollCtl", zap.Error(err))
					} else if ptr.newConnectHandler != nil {
						ptr.newConnectHandler(connFd)
					}
				} else {
					go func(fd int) { //c处理
						for {
							dataBuf := C.newMsg(C.int(fd))
							if dataBuf == nil {
								if ptr.closeConnectHandler != nil {
									ptr.closeConnectHandler(fd)
								}
								if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, fd, &event); err != nil {
									logger.Error("tcp close", zap.Error(err))
								}
								syscall.Close(fd)
								break
							} else if ptr.newMessageHandler != nil && len(C.GoString((*C.char)(dataBuf))) > 0 {
								ptr.newMessageHandler(fd, []byte(C.GoString((*C.char)(dataBuf))))
							}
							C.destruct1((**C.char)(&dataBuf))
						}

					}(int(events[ev].Fd))
				}
			}
		}
	}).Catch(func(err interface{}) {
		logger.ErrorString("tcp", "run error", fmt.Sprintf(err))
		os.Exit(1)
	}).Run()
}
