package entity

import (
	"context"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/util/structure"
	"gorm.io/gorm"
)

// GetRoleRouterDB 获取用户角色关联储存
func GetRoleRouterDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(RoleRouter))
}

// SchemaRoleRouter 角色路由资源
type SchemaRoleRouter schema.RoleRouter

// ToRoleRouter转换为角色路由资源实体
func (a SchemaRoleRouter) ToRoleRouter() *RoleRouter {
	item := new(RoleRouter)
	structure.Copy(a, item)
	return item
}

// RoleRouter 角色路由关联实体
type RoleRouter struct {
	ID       string `gorm:"column:id;primaryKey;size:36;"`
	RoleID   string `gorm:"column:role_id;size:36;index;default:'';not null;"`   // 角色内码
	RouterID string `gorm:"column:router_id;size:36;index;default:'';not null;"` // 路由资源内码
}

// ToSchemaRoleRouter 转换为角色路由资源对象
func (a RoleRouter) ToSchemaRoleRouter() *schema.RoleRouter {
	item := new(schema.RoleRouter)
	structure.Copy(a, item)
	return item
}

// 角色关联路由资源列表
type RoleRouters []*RoleRouter

// ToSchemaRoleRouters 转换为角色路由对象列表
func (a RoleRouters) ToSchemaRoleRouters() []*schema.RoleRouter {
	list := make([]*schema.RoleRouter, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleRouter()
	}

	return list
}
