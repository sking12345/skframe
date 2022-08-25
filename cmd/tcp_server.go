package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"skframe/cmd/tcp"
	"skframe/pkg/config"
	"skframe/pkg/logger"
	"skframe/pkg/try"
)

var TcpServer = &cobra.Command{
	Use:   "tcp",
	Short: "Start websocket server",
	Run:   runTcpServer,
	Args:  cobra.NoArgs,
}

func runTcpServer(cmd *cobra.Command, args []string) {
	try.NewTry(func() {
		port := config.GetInt("tcp.port")
		if port <= 0 {
			panic("not set tcp listen port")
		}
		ptr := tcp.NewTcp(port)
		ptr.SetNewConnectHandler(func(fd int) {
			handler := config.GetInterface("tcp.newConnect")
			if handler != nil {
				handler.(func(fd int))(fd)
			}
		})
		ptr.SetCloseConnectHandler(func(fd int) {
			handler := config.GetInterface("tcp.closeConnect")
			if handler != nil {
				handler.(func(fd int))(fd)
			}
		})
		ptr.SetNewMessageHandler(func(fd int, data []byte) {
			handler := config.GetInterface("tcp.newMessage")
			if handler != nil {
				handler.(func(fd int))(fd)
			}
		})
		ptr.Run()
	}).Catch(func(err interface{}) {
		logger.ErrorString("tcp", "run fail:", fmt.Sprintf("%s", err))
		os.Exit(1)
	}).Run()
}
