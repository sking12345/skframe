package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"skframe/pkg/console"
	"skframe/pkg/helpers"
	"skframe/pkg/jwt"
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

func Cors() gin.HandlerFunc { //跨域
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		// 处理请求
		c.Next()
	}
}

func RegisterAPIRoutes(router *gin.Engine) {
	router.Use(Cors())
	router.GET("/test", func(context *gin.Context) {
		fmt.Println("xx")
	})

}
