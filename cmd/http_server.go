package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"skframe/bootstrap"
	"skframe/pkg/config"
	"skframe/pkg/console"
	"skframe/pkg/logger"
)

//CmdHttpServer

var HttpServer = &cobra.Command{
	Use:   "http",
	Short: "Start web server",
	Run:   runHttpWeb,
	Args:  cobra.NoArgs,
}

func runHttpWeb(cmd *cobra.Command, args []string) {
	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	bootstrap.SetupRoute(router)
	console.Info("start http,port:" + config.Get("app.port"))
	if err := router.Run(":" + config.Get("app.port")); err != nil {
		logger.ErrorString("CMD", "httpServer", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
}
