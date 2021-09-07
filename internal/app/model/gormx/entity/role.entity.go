package entity

import (
	"context"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/util/structure"
	"gorm.io/gorm"
	"time"
)

// GetRoleDB 获取角色储存
func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Role))
}

// SchemaRole 角色对象
type SchemaRole schema.Role

// ToRole 转换为角色实体
func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	structure.Copy(a, item)
	return item
}

// Role 角色实体
type Role struct {
	ID        string         `gorm:"column:id;primaryKey;size:36;"`
	Name      string         `gorm:"column:name;size:100;index;default:'';not nul;"` // 角色名称
	Sequence  int            `gorm:"column:sequence;index;default:0;not null;"`      // 排序值
	Memo      *string        `gorm:"column:memo;size:1024;"`                         // 备注
	Status    int            `gorm:"column:status;index;default:0;not null;"`        // 状态(1:启用 2:禁用)
	Creator   string         `gorm:"column:creator;size:36;"`                        // 创建者
	CreatedAt time.Time      `gorm:"column:created_at;index;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;"`
}

// ToSchemaRole 转换角色对象
func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	structure.Copy(a, item)
	return item
}

// Roles 角色实体列表
type Roles []*Role

// ToSchemaRoles 转换角色对象列表
func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}
