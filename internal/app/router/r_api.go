package router

import (
	_ "PeachyTalkService/docs"
	"PeachyTalkService/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAPI 注册 api 组路由器
func (a *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")

	g.Use(middleware.UserAuthMiddleware(a.Auth,
		middleware.AllowPathPrefixSkipper("/api/v1/pub/login"),
	))

	g.Use(middleware.RateLimiterMiddleware())

	v1 := g.Group("/v1")
	{
		pub := v1.Group("/pub")
		{
			gLogin := pub.Group("login")
			{
				gLogin.GET("captchaid", a.LoginAPI.GetCaptcha)
				gLogin.GET("captcha", a.LoginAPI.ResCaptcha)
				gLogin.POST("", a.LoginAPI.Login)
				gLogin.POST("exit", a.LoginAPI.Logout)
			}

			gCurrent := pub.Group("current")
			{
				gCurrent.PUT("password", a.LoginAPI.UpdatePassword)
				gCurrent.GET("user", a.LoginAPI.GetUserInfo)
			}
			pub.POST("/refresh-token", a.LoginAPI.RefreshToken)
		}

		management := v1.Group("/management")
		{
			gDemo := management.Group("demos")
			{
				gDemo.GET("", a.DemoAPI.Query)
				gDemo.GET(":id", a.DemoAPI.Get)
				gDemo.POST("", a.DemoAPI.Create)
				gDemo.PUT(":id", a.DemoAPI.Update)
				gDemo.DELETE(":id", a.DemoAPI.Delete)
				gDemo.PATCH(":id/enable", a.DemoAPI.Enable)
				gDemo.PATCH(":id/disable", a.DemoAPI.Disable)
			}

			gUser := management.Group("users")
			{
				gUser.GET("", a.UserAPI.Query)
				gUser.GET(":id", a.UserAPI.Get)
				gUser.POST("", a.UserAPI.Create)
				gUser.PUT(":id", a.UserAPI.Update)
				gUser.DELETE(":id", a.UserAPI.Delete)
				gUser.PATCH(":id/enable", a.UserAPI.Enable)
				gUser.PATCH(":id/disable", a.UserAPI.Disable)
			}
		}

	}
}
