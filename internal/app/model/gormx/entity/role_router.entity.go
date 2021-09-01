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
type SchemaRoleRouter schema.RouterResource

// ToRoleRouter转换为角色路由资源实体
func (a SchemaRoleRouter) ToSchemaRoleRouter() *RoleRouter {
	item := new(RoleRouter)
	structure.Copy(a, item)
	return item
}

// RoleRouter 角色路由关联实体
type RoleRouter struct {
	ID       string `gorm:"column:id;primaryKey;size:36;"`
	RoleID   string `gorm:"column:user_id;size:36;index;default:'';not null;"`   // 角色内码
	RouterID string `gorm:"column:router_id;size:36;index;default:'';not null;"` // 路由资源内码
}

// ToSchemaRoleRouter 转换为角色路由资源对象
func (a UserRole) ToSchemaRoleRouter() *schema.RouterResource {
	item := new(schema.RouterResource)
	structure.Copy(a, item)
	return item
}
