package config

import (
	"fmt"
	"skframe/pkg/config"
)

func init() {
	config.Add("udp", func() map[string]interface{} {
		return map[string]interface{}{
			"port":     config.Env("UDP_PORT", 3802),
			"buffSize": config.Env("UDP_BUFF_SIZE", 1024),
			"msgHandler": func(fd int, data []byte, addr []byte) {
				fmt.Println(data)
			},
		}
	})
}
