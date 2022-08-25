package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"skframe/cmd/udp"
	"skframe/pkg/config"
	"skframe/pkg/logger"
	"skframe/pkg/try"
)

var UdpServer = &cobra.Command{
	Use:   "udp",
	Short: "Start websocket server",
	Run:   runUdpServer,
	Args:  cobra.NoArgs,
}

func runUdpServer(cmd *cobra.Command, args []string) {
	try.NewTry(func() {
		port := config.GetInt("udp.port")
		if port <= 0 {
			panic("not set tcp listen port")
		}
		ptr := udp.NewServer(port, config.GetInt("udp.buffSize"))
		ptr.SetMessageHandler(func(fd int, data []byte, addr []byte) {
			handler := config.GetInterface("udp.msgHandler")
			handler.(func(fd int, data []byte, addr []byte))(fd, data, addr)
		})
		ptr.Run()
	}).Catch(func(err interface{}) {
		logger.ErrorString("udp", "run fail:", fmt.Sprintf("%s", err))
		os.Exit(1)
	}).Run()
}
