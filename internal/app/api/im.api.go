package api

import (
	"ginAdmin/internal/app/ginx"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"ginAdmin/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

// IMSet 注入IM
var IMSet = wire.NewSet(wire.Struct(new(IM), "*"))

// IM 聊天管理
type IM struct {
	IMSrv *service.IM
}

// Auth Login IM 账号认证
// @Tags IM
// @Summary IM
// @Param body body schema.IMClient true "请求参数"
// @Success 200
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/im/auth [post]
func (a *IM) Auth(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.IMClient
	if err := ginx.ParseJSON(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.IMSrv.Verify(ctx, params.Username, params.Password)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ctx = logger.NewUserIDContext(ctx, params.Username)
	ctx = logger.NewTagContext(ctx, "__login__")
	logger.WithContext(ctx).Infof("登入系统IM")
	c.Status(http.StatusOK)
}

// Superuser Login IM 超级账号认证
// @Tags IM
// @Summary IM
// @Param body schema.IMClient true "请求参数"
// @Success 200
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/im/auth [post]
func (a *IM) Superuser(c *gin.Context) {
	var params schema.IMClient
	if err := ginx.ParseJSON(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.IMSrv.IsSuperUser(params.Username)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Acl Login IM 权限认证
// @Tags IM
// @Summary IM
// @Param body schema.IMAcl true "请求参数"
// @Success 200
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/im/acl [post]
func (a *IM) Acl(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.IMAcl
	if err := ginx.ParseJSON(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.IMSrv.AclVerify(ctx, params.Username, params.Topic, params.Access)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
