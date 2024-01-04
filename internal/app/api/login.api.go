package api

import (
	"PeachyTalkService/internal/app/config"
	"PeachyTalkService/internal/app/ginx"
	"PeachyTalkService/internal/app/schema"
	"PeachyTalkService/internal/app/service"
	"PeachyTalkService/pkg/errors"
	"PeachyTalkService/pkg/logger"

	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// LoginSet 注入Login
var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

// Login 登录管理
type Login struct {
	LoginSrv *service.Login
}

// GetCaptcha 获取验证码信息
// @Tags 登录管理
// @Summary 获取验证码信息
// @Success 200 {object} schema.LoginCaptcha
// @Router /api/v1/pub/login/captchaid [get]
func (a *Login) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginSrv.GetCaptcha(ctx, config.C.Captcha.Length)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// ResCaptcha 响应图形验证码
// @Tags 登录管理
// @Summary 响应图形验证码
// @Param id query string true "验证码ID"
// @Param reload query string false "重新加载"
// @Produce image/png
// @Success 200 "图形验证码"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/login/captcha [get]
func (a *Login) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		ginx.ResError(c, errors.New400Response("请提供验证码ID"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginx.ResError(c, errors.New400Response("未找到验证码ID"))
			return
		}
	}

	cfg := config.C.Captcha
	err := a.LoginSrv.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginx.ResError(c, err)
	}
}

// Login 用户登录
// @Tags 登录管理
// @Summary 用户登录
// @Param body body schema.LoginParam true "请求参数"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/login [post]
func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginx.ResError(c, errors.New400Response("无效的验证码"))
		return
	}

	user, err := a.LoginSrv.Verify(ctx, item.UserName, item.Password)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	userID := user.ID
	// 将用户ID放入上下文
	ginx.SetUserID(c, userID)

	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, userID)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, userID)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("登入系统")
	ginx.ResSuccess(c, tokenInfo)
}

// Logout 用户登出
// @Tags 登录管理
// @Summary 用户登出
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Router /api/v1/pub/login/exit [post]
func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	// 检查用户是否处于登录状态，如果是则执行销毁
	userID := ginx.GetUserID(c)
	if userID != "" {
		ctx = logger.NewTagContext(ctx, "__logout__")
		err := a.LoginSrv.DestroyToken(ctx, ginx.GetToken(c))
		if err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
		logger.WithContext(ctx).Infof("登出系统")
	}
	ginx.ResOK(c)
}

// RefreshToken 刷新令牌
// @Tags 登录管理
// @Summary 刷新令牌
// @Security ApiKeyAuth
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/refresh-token [post]
func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, tokenInfo)
}

// GetUserInfo 获取当前用户信息
// @Tags 登录管理
// @Summary 获取当前用户信息
// @Security ApiKeyAuth
// @Success 200 {object} schema.UserLoginInfo
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/user [get]
func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginSrv.GetLoginInfo(ctx, ginx.GetUserID(c))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, info)
}

// UpdatePassword 更新个人密码
// @Tags 登录管理
// @Summary 更新个人密码
// @Security ApiKeyAuth
// @Param body body schema.UpdatePasswordParam true "请求参数"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/pub/current/password [put]
func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.LoginSrv.UpdatePassword(ctx, ginx.GetUserID(c), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
