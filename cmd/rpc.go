package cmd

import (
	"github.com/spf13/cobra"
)

var RPCServer = &cobra.Command{
	Use:   "http",
	Short: "Start web server",
	Run:   runRPC,
	Args:  cobra.NoArgs,
}

func runRPC(cmd *cobra.Command, args []string) {
	//rpc := rpc.Micro{}
	//handler := config.GetInterface("micro.handler")
	//if handler == nil {
	//	console.Exit("rpc.micro handler nil")
	//}
	//err := rpc.Start(
	//	config.GetString("micro.addr"),
	//	config.GetString("micro.name"),
	//	config.GetString("micro.version"),
	//	config.GetUint("micro.ttl"),
	//	config.GetUint("micro.interval"),
	//	handler.(func(service micro.Service)),
	//	)
	//if err != nil {
	//	console.Exit("Unable to start rpc, error:" + err.Error())
	//}
}

