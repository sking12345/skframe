package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"skframe/app/controlles"
	"skframe/pkg/console"
	"skframe/pkg/helpers"
	"skframe/pkg/jwt"
	"skframe/pkg/ws"
)

func authJWT() gin.HandlerFunc { //用户token 验证
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if helpers.Empty(authHeader) == true {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。",
			})
		}
		claims, err := jwt.NewJWT().ParserToken(authHeader)
		if err != nil {
			console.Error(err.Error())
		}
		console.Info(fmt.Sprintf("userId:%d", claims.UserId))
		ctx.Set("userId", claims.UserId)
		ctx.Next()
	}
}
func guestJWT() gin.HandlerFunc { //游客访问
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func RegisterAPIRoutes(ginEngine *gin.Engine) {
	userCtl := new(controlles.UserController)
	ginEngine.POST("/register", guestJWT(), func(context *gin.Context) {
		userCtl.Register(context)
	})
	ginEngine.POST("/login", guestJWT(), func(context *gin.Context) {
		userCtl.Login(context)
	})
	ginEngine.POST("/add/friend", guestJWT(), func(context *gin.Context) {
		userCtl.AddFriend(context)
	})

}
func wsJwt(ctx *ws.Context) {
	ctx.Next()
}

func RegisterSocketRoutes(ctx *ws.Engine) {
	//testCtl := new(controlles.TestController)
	//ctx.GET("/", wsJwt, func(context *ws.Context) {
	//	testCtl.CreatTokenTest(context)
	//})
}
