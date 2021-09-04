package api

import (
	"ginAdmin/internal/app/ginx"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// RouterResourceSet 注入RouterResource
var RouterResourceSet = wire.NewSet(wire.Struct(new(RouterResource), "*"))

// RouterResource 路由资源
type RouterResource struct {
	RouterResourceSrv *service.RouterResource
}

// Query 查询数据
// @Tags 路由资源
// @Security ApiKeyAuth
// @Summary 查询数据
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Param roleId query string false "角色ID"
// @Param excludeIDs query []string false "需要排除的ID"
// @Param status query int 0 "启用状态"
// @Success 200 {object} schema.ListResult{list=[]schema.RouterResource} "查询结果"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/management/routerResources [get]
func (a *RouterResource) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RouterResourceQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.RouterResourceSrv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
// @Tags 路由资源
// @Security ApiKeyAuth
// @Summary 查询指定数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.RouterResource
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/management/routerResources/{id} [get]
func (a *RouterResource) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RouterResourceSrv.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// Create 创建数据
// @Tags 路由资源
// @Security ApiKeyAuth
// @Summary 创建数据
// @Param body body schema.RouterResource true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/management/routerResources [post]
func (a *RouterResource) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.RouterResource
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}
	item.Creator = ginx.GetUserID(c)
	result, err := a.RouterResourceSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// Update 更新数据
// @Tags 路由资源
// @Security ApiKeyAuth
// @Summary 更新数据
// @Param id path string true "唯一标识"
// @Param body body schema.RouterResource true "更新数据"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/management/routerResources/{id} [put]
func (a *RouterResource) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.RouterResource
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.RouterResourceSrv.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResOK(c)
}

// Delete 删除数据
// @Tags 路由资源
// @Security ApiKeyAuth
// @Summary 删除数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/management/routerResources/{id} [delete]
func (a *RouterResource) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RouterResourceSrv.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Enable 启用数据
// @Tags 路由资源
// @Security ApiKeyAuth
// @Summary 启用数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/management/routerResources/{id}/enable [patch]
func (a *RouterResource) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RouterResourceSrv.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Disable 禁用数据
// @Tags 路由资源
// @Security ApiKeyAuth
// @Summary 禁用数据
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/management/routerResources/{id}/disable [patch]
func (a *RouterResource) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RouterResourceSrv.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
