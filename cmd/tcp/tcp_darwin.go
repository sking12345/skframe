//+build darwin

package tcp

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
	if (readSizze < 0){
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

////类似epoll_create
//int kqueue(void);
////组合了epoll_ctl及epoll_wait功能，changelist与nchanges为修改列表，eventlist与nevents为返回的事件列表
//int kevent(int kq, const struct kevent *changelist, int nchanges, struct kevent *eventlist, int nevents, const struct timespec *timeout);
////设定kevent参数的宏，详细解释见后面的调用示例
//EV_SET(&kev, ident, filter, flags, fflags, data, udata);

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
		//var event syscall.EpollEvent
		socketfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			panic(err)
		}
		defer syscall.Close(socketfd)
		if err = syscall.SetNonblock(socketfd, true); err != nil {
			panic(err)
		}
		addr := syscall.SockaddrInet4{Port: ptr.Port}
		copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
		syscall.Bind(socketfd, &addr)
		syscall.Listen(socketfd, 10) //backlog:同时处理最大连接数
		changeEvent := syscall.Kevent_t{
			Ident:  uint64(socketfd),
			Filter: syscall.EVFILT_READ,
			Flags:  syscall.EV_ADD | syscall.EV_CLEAR, //EV_CLEAR : 用户检索事件后，将重置其状态。否则会不停触发事件
			Fflags: 0,
			Data:   0,
			Udata:  nil,
		}
		kQueue, err := syscall.Kqueue()
		if err != nil {
			panic(err)
		}
		changeEventRegistered, err := syscall.Kevent(
			kQueue,
			[]syscall.Kevent_t{changeEvent},
			nil,
			nil,
		)
		if err != nil {
			panic(err)
		}
		if changeEventRegistered == -1 {
			panic("syscall.Kevent error(=-1)")
		}
		logger.Info("tcp", zap.Any("status", "success"), zap.Any("port", ptr.Port))
		for {
			newEvents := make([]syscall.Kevent_t, 10)
			numNewEvents, err := syscall.Kevent(
				kQueue,
				nil,
				newEvents,
				nil,
			)
			if err != nil {
				continue
			}
			for i := 0; i < numNewEvents; i++ {
				currentEvent := newEvents[i]
				eventFileDescriptor := int(currentEvent.Ident)
				if currentEvent.Flags&syscall.EV_EOF != 0 {
					if ptr.closeConnectHandler != nil {
						ptr.closeConnectHandler(eventFileDescriptor)
					}
					socketEvent := syscall.Kevent_t{
						Ident:  uint64(eventFileDescriptor),
						Filter: syscall.EVFILT_READ,
						Flags:  syscall.EV_DELETE | syscall.EV_CLEAR,
						Fflags: 0,
						Data:   0,
						Udata:  nil,
					}
					_, err := syscall.Kevent(
						kQueue,
						[]syscall.Kevent_t{socketEvent},
						nil,
						nil,
					)
					if err != nil {
						panic(err)
					}
					syscall.Close(eventFileDescriptor)
				} else if eventFileDescriptor == socketfd {
					socketConnection, _, err := syscall.Accept(eventFileDescriptor)
					if err != nil {
						continue
					}
					//EV_CLEAR : 用户检索事件后，将重置其状态。否则会不停触发事件
					socketEvent := syscall.Kevent_t{
						Ident:  uint64(socketConnection),
						Filter: syscall.EVFILT_READ,
						Flags:  syscall.EV_ADD | syscall.EV_CLEAR,
						Fflags: 0,
						Data:   0,
						Udata:  nil,
					}
					socketEventRegistered, err := syscall.Kevent(
						kQueue,
						[]syscall.Kevent_t{socketEvent},
						nil,
						nil,
					)
					if err != nil || socketEventRegistered == -1 {
						continue
					} else if ptr.newConnectHandler != nil {
						ptr.newConnectHandler(eventFileDescriptor)
					}
				} else if currentEvent.Filter&syscall.EVFILT_READ != 0 {
					go func(fd int) {
						for { //当多个少于mute的数据同时发送到服务端时,事件被重置，不会触发
							dataBuf := C.newMsg(C.int(fd))
							if dataBuf == nil {
								break
							}
							if len(C.GoString((*C.char)(dataBuf))) <= 0 {
								break
							}
							if ptr.newMessageHandler != nil && len(C.GoString((*C.char)(dataBuf))) > 0 {
								ptr.newMessageHandler(fd, []byte(C.GoString((*C.char)(dataBuf))))
							}
							C.destruct1((**C.char)(&dataBuf))
						}
					}(eventFileDescriptor)
				}
			}
		}
	}).Catch(func(err interface{}) {
		logger.ErrorString("tcp", "run error", fmt.Sprintf("%s", err))
		os.Exit(1)
	}).Run()
}
