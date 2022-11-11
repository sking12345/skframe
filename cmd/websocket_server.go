package cmd

import (
	"github.com/spf13/cobra"
)

//CmdHttpServer

var WebSocketServer = &cobra.Command{
	Use:   "websocket",
	Short: "Start websocket server",
	Run:   runWebSocketServer,
	Args:  cobra.NoArgs,
}

func runWebSocketServer(cmd *cobra.Command, args []string) {
	//gin.SetMode(gin.ReleaseMode)
	//req := gin.Default()
	//wsCtx := ws.Engine{}
	//
	//path := config.GetString("ws.path")
	//req.Any("/"+path, func(ctx *gin.Context) {
	//	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	//		return true
	//	}}
	//	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	//	defer conn.Close()
	//	if err != nil {
	//		logger.Error("websocket server", zap.Error(err))
	//		return
	//	}
	//	clientCtx := wsCtx.NewConnect(conn)
	//	for {
	//		_, data, err := conn.ReadMessage()
	//		if err != nil {
	//			logger.Info("websocket server", zap.Error(err))
	//			wsCtx.CloseConnect(clientCtx)
	//			return
	//		}
	//		wsCtx.NewMessage(data, clientCtx)
	//	}
	//})
	//req.Run(":" + config.GetString("ws.port"))

}
