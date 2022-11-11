package config

import (
	"fmt"
	"skframe/cmd/udp"
	"skframe/pkg/config"
	"unsafe"
)

type udpPackBuf struct {
	ReceiveUid int32
	SendUid int32
	UniqueNum int32
	Index int32
	TotalSize int32
	DataBuf []int8
}




func init() {
	config.Add("udp", func() map[string]interface{} {
		return map[string]interface{}{
			"port":     config.Env("UDP_PORT", 3802),
			"buffSize": config.Env("UDP_BUFF_SIZE", 1024),
			"msgHandler": func(fd int, data []byte, addr []byte) {
				var ptestStruct  = *(**udpPackBuf)(unsafe.Pointer(&data))

				fmt.Println(ptestStruct.SendUid)
				fmt.Println(ptestStruct.ReceiveUid)
				fmt.Println(ptestStruct.TotalSize)
				fmt.Println(ptestStruct.DataBuf)
				fmt.Println(string(data))
				udp.SendData(fd,data,addr)

			},
		}
	})
}
