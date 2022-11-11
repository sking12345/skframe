package config

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"skframe/app/middlewares"
	"skframe/pkg/config"
	"skframe/routes"
	"strings"
)

func init()  {
	config.Add("web", func() map[string]interface{} {

		return map[string]interface{}{
			"port":		config.Env("APP_PORT","3000"),
			"key":		config.Env("APP_KEY","123456"),
			// 用以生成链接
			"url": 		config.Env("APP_URL", "http://localhost:3000"),
			// 设置时区，JWT 里会使用，日志记录里也会使用到
			"timezone": config.Env("TIMEZONE", "Asia/Shanghai"),
			// API 域名，未设置的话所有 API URL 加 api 前缀，如 http://domain.com/api/v1/users
			"api_domain": config.Env("API_DOMAIN"),
			"routeHandler": func(router *gin.Engine) {
				router.Use(
					middlewares.Logger(),
					middlewares.Recovery(),
				)
				routes.RegisterAPIRoutes(router)
				router.NoRoute(func(c *gin.Context) {
					// 获取标头信息的 Accept 信息
					acceptString := c.Request.Header.Get("Accept")
					if strings.Contains(acceptString, "text/html") {
						// 如果是 HTML 的话
						c.String(http.StatusNotFound, "页面返回 404")
					} else {
						// 默认返回 JSON
						c.JSON(http.StatusNotFound, gin.H{
							"error_code":    404,
							"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
						})
					}
				})
			},
		}
	})
}