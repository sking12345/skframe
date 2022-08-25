package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"skframe/bootstrap"
	"skframe/cmd"
	cmdMake "skframe/cmd/make"
	"skframe/config"
	pkCofing "skframe/pkg/config"
	"skframe/pkg/console"
)

func init() {
	config.Initialize()
}

func main() {

	var rootCmd = &cobra.Command{
		Use:   "skFrame",
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,
		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			pkCofing.InitConfig(cmd.Env)
			//bootstrap.SetupCache()
			bootstrap.SetupLogger()
			//bootstrap.SetupDB()
			//bootstrap.SetUpEtcd()
		},
	}
	//注册子命令
	rootCmd.AddCommand(cmd.HttpServer, cmd.WebSocketServer, cmdMake.CmdMake, cmd.TcpServer, cmd.UdpServer,cmd.RPCServer)
	cmd.RegisterDefaultCmd(rootCmd, cmd.HttpServer)
	cmd.RegisterGlobalFlags(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
