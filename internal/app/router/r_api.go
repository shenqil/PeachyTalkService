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
		middleware.AllowPathPrefixSkipper("/api/v1/login"),
	))

	g.Use(middleware.RateLimiterMiddleware())

	v1 := g.Group("/v1")
	{
		gLogin := v1.Group("login")
		{
			gLogin.GET("captchaid", a.LoginAPI.GetCaptcha)
			gLogin.GET("captcha", a.LoginAPI.ResCaptcha)
			gLogin.POST("", a.LoginAPI.Login)
			gLogin.POST("exit", a.LoginAPI.Logout)
		}

		gCurrent := v1.Group("current")
		{
			gCurrent.PUT("password", a.LoginAPI.UpdatePassword)
			gCurrent.GET("user", a.LoginAPI.GetUserInfo)
		}
		v1.POST("/refresh-token", a.LoginAPI.RefreshToken)

		gDemo := v1.Group("demos")
		{
			gDemo.GET("", a.DemoAPI.Query)
			gDemo.GET(":id", a.DemoAPI.Get)
			gDemo.POST("", a.DemoAPI.Create)
			gDemo.PUT(":id", a.DemoAPI.Update)
			gDemo.DELETE(":id", a.DemoAPI.Delete)
			gDemo.PATCH(":id/enable", a.DemoAPI.Enable)
			gDemo.PATCH(":id/disable", a.DemoAPI.Disable)
		}

		gUser := v1.Group("users")
		{
			gUser.GET("", a.UserAPI.Query)
			gUser.GET(":id", a.UserAPI.Get)
			gUser.POST("", a.UserAPI.Create)
			gUser.PUT(":id", a.UserAPI.Update)
			gUser.DELETE("", a.UserAPI.BatchDelete)
			gUser.DELETE(":id", a.UserAPI.Delete)
			gUser.PATCH(":id/enable", a.UserAPI.Enable)
			gUser.PATCH(":id/disable", a.UserAPI.Disable)
		}

		gFile := v1.Group("file")
		{
			gFile.POST("", a.File.Upload)
		}
	}
}
