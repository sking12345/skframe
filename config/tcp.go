package config

import (
	"skframe/pkg/config"
	"skframe/routes"
)

func init() {
	config.Add("tcp", func() map[string]interface{} {
		return map[string]interface{}{
			"port":         config.Env("TCP_PORT", 3800),
			"newConnect":   routes.TCPNewConnect,
			"closeConnect": routes.TCPCloseConnect,
			"newMessage":   routes.TCPNewMessage,
		}
	})
}
