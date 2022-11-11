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
	"skframe/pkg/database"
)

func init() {
	config.Initialize()
}

func main() {
	pkCofing.InitConfig(cmd.Env)
	bootstrap.SetupLogger()
	bootstrap.SetupDB()
	count := database.Count("test", map[string]interface{}{"id": map[string]interface{}{">": 1}}, nil)
	fmt.Println(count)
	//tableFieldInfo := map[string]string{
	//	"id":      "uint64",
	//	"name":    "string",
	//	"number":  "int",
	//	"num1":    "decimal",
	//	"testcol": "float",
	//}
	//_, reslut := database.Find("test", tableFieldInfo, map[string]interface{}{"id": 2}, nil)
	//fmt.Println(reslut)
	//err, _ := database.Create("test", map[string]interface{}{"name": "dddd", "number": "xx"})
	//fmt.Println(err)
	//database.Update("test", map[string]interface{}{"name": "xxxdd", "num1": 123}, map[string]interface{}{
	//	"id": map[string]interface{}{">": 1}},
	//)
	//connectInfo, err := database.Begin()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//id, _ := database.Del("test", map[string]interface{}{"id": 4}, connectInfo.Tx)
	//database.Commit(connectInfo)
	//fmt.Println(id)
	//bootstrap.DestructDB()
	//fmt.Println(reslut)
	return

	var rootCmd = &cobra.Command{
		Use:   "skFrame",
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,
		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			pkCofing.InitConfig(cmd.Env)
			//bootstrap.SetupCache()
			bootstrap.SetupLogger()
			bootstrap.SetupDB()
			//bootstrap.SetUpEtcd()
		},
	}
	//注册子命令
	rootCmd.AddCommand(cmd.HttpServer,
		cmd.WebSocketServer,
		cmdMake.CmdMake,
		cmd.TcpServer,
		cmd.UdpServer,
		//cmd.RPCServer,
	)
	cmd.RegisterDefaultCmd(rootCmd, cmd.HttpServer)
	//cmd.RegisterDefaultCmd(rootCmd, cmd.UdpServer)
	cmd.RegisterGlobalFlags(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
