package entity

import (
	"context"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/util/structure"
	"gorm.io/gorm"
	"time"
)

// GetRouterResourceDB 获取路由资源
func GetRouterResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(RouterResource))
}

// SchemaRouterResource 路由资源对象
type SchemaRouterResource schema.RouterResource

// ToRouterResource 转换为角色实体
func (a SchemaRouterResource) ToRouterResource() *RouterResource {
	item := new(RouterResource)
	structure.Copy(a, item)
	return item
}

// RouterResource 路由相关资源
type RouterResource struct {
	ID        string         `gorm:"column:id;primary_key;size:36;"`
	Name      string         `gorm:"column:name;size:100;default:'';not null;"`
	Memo      string         `gorm:"column:memo;size:1024;default:'';not null;"`
	Method    string         `gorm:"column:method;size:100;default:'';not null;"` // 资源请求方式(支持正则)
	Path      string         `gorm:"column:path;size:100;default:'';not null;"`   // 资源请求路径（支持/:id匹配）
	Status    int            `gorm:"column:status;index;default:0;not null;"`
	Creator   string         `gorm:"column:creator;size:36;"` // 创建人
	CreatedAt time.Time      `gorm:"column:created_at;index;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;"`
}

// ToSchemaRouterResource 转换为角色对象
func (a RouterResource) ToSchemaRouterResource() *schema.RouterResource {
	item := new(schema.RouterResource)
	structure.Copy(a, item)
	return item
}

// RouterResources 路由资源列表
type RouterResources []*RouterResource

// ToSchemaRouterResources 转换为路由资源列表
func (a RouterResources) ToSchemaRouterResources() []*schema.RouterResource {
	list := make([]*schema.RouterResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRouterResource()
	}
	return list
}
