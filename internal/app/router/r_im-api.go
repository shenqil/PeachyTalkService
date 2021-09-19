package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Register IM-API 注册 IM-API 组路由器
func (a *Router) RegisterIMAPI(app *gin.Engine) {
	g := app.Group("/api")

	v1 := g.Group("/v1")
	{
		im := v1.Group("/im")
		{
			// 测试服务是否注册成功
			im.GET("ping", func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			})

			im.POST("auth", a.IMApi.Auth)
			im.POST("acl", a.IMApi.Acl)
			im.POST("superuser", a.IMApi.Superuser)
		}
	}
}
