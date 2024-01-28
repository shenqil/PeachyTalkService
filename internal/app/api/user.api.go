package api

import (
	"PeachyTalkService/internal/app/ginx"
	"PeachyTalkService/internal/app/schema"
	"PeachyTalkService/internal/app/service"
	"PeachyTalkService/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

// User 用户管理
type User struct {
	UserSrv *service.User
}

// Query 查询数据
// @Tags 用户管理
// @Summary 查询数据
// @Security ApiKeyAuth
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Param status query int false "状态(1:启用 2:停用)"
// @Success 200 {object} schema.ListResult{list=[]schema.Users} "查询结果"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.UserSrv.Query(ctx, params)
	result.Data.CleanSecure()
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
// @Tags 用户管理
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.User
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserSrv.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item.CleanSecure())
}

// Create 创建数据
// @Tags 用户管理
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.User true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	} else if item.Password == "" {
		ginx.ResError(c, errors.New400Response("密码不能为空"))
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.UserSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// Update 更新数据
// @Tags 用户管理
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Param body body schema.User true "更新数据"
// @Success 200 {object} schema.User
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.UserSrv.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Delete 删除数据
// @Tags 用户管理
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSrv.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// BatchDelete 批量删除数据
// @Tags 用户管理
// @Summary 批量删除数据
// @Security ApiKeyAuth
// @Param body body schema.IDsRequest true "需要删除的"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users [delete]
func (a *User) BatchDelete(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.IDsRequest
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.UserSrv.BatchDelete(ctx, item.IDs)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Enable 启用数据
// @Tags 用户管理
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/enable [patch]
func (a *User) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSrv.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Disable 禁用数据
// @Tags 用户管理
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/users/{id}/disable [patch]
func (a *User) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserSrv.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
