package api

import (
	"ginAdmin/internal/app/ginx"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// MenuSet 注入Menu
var MenuSet = wire.NewSet(wire.Struct(new(Menu), "*"))

// Menu 菜单管理
type Menu struct {
	MenuSrv *service.Menu
}

// Query 查询数据
// @Tags 菜单管理
// @Summary 查询数据
// @Security ApiKeyAuth
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Param showStatus query int false "显示状态(1:显示 2:隐藏)"
// @Param parentID query string false "父级ID"
// @Success 200 {object} schema.ListResult{list=[]schema.Menu} "查询结果"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus [get]
func (a *Menu) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.MenuSrv.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

// QueryTree 查询菜单树
// @Tags 菜单管理
// @Summary 查询菜单树
// @Security ApiKeyAuth
// @Param status query int false "状态(1:启用 2:禁用)"
// @Param parentID query string false "父级ID"
// @Success 200 {object} schema.ListResult{list=[]schema.MenuTree} "查询结果"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus.tree [get]
func (a *Menu) QueryTree(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.MenuQueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.MenuSrv.Query(ctx, params, schema.MenuQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResList(c, result.Data.ToTree())
}

// Get 查询指定数据
// @Tags 菜单管理
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.Menu
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id} [get]
func (a *Menu) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.MenuSrv.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// Create 创建数据
// @Tags 菜单管理
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.Menu true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus [post]
func (a *Menu) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.MenuSrv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// Update 更新数据
// @Tags 菜单管理
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Param body body schema.Menu true "更新数据"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id} [put]
func (a *Menu) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Menu
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.MenuSrv.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Delete 删除数据
// @Tags 菜单管理
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id} [delete]
func (a *Menu) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Enable 启用数据
// @Tags 菜单管理
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id}/enable [patch]
func (a *Menu) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Disable 禁用数据
// @Tags 菜单管理
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path string true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/menus/{id}/disable [patch]
func (a *Menu) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.MenuSrv.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
