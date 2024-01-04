package entity

import (
	"PeachyTalkService/internal/app/schema"
	"PeachyTalkService/pkg/util/structure"
	"context"
	"time"

	"gorm.io/gorm"
)

// GetGroupDB 获取群组数据
func GetGroupDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Group))
}

// SchemaGroup 群组对象
type SchemaGroup schema.Group

// ToGroup 转换为群组实体
func (a SchemaGroup) ToGroup() *Group {
	item := new(Group)
	structure.Copy(a, item)
	return item
}

// Group 群组实体
type Group struct {
	ID        string         `gorm:"column:id;primaryKey;size:36;"`                        // ID
	GroupName string         `gorm:"column:group_name;size:64;index;default:'';not null;"` // 群名称
	Brief     string         `gorm:"column:brief;size:1024;default:'';not null;"`          // 简介
	Avatar    string         `gorm:"column:avatar;size:128;default:'';not null;"`          // 头像
	Owner     string         `gorm:"column:owner;size:36;"`                                // 拥有者
	Creator   string         `gorm:"column:creator;size:36;"`                              // 创建者
	CreatedAt time.Time      `gorm:"column:created_at;index;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;"`
}

// ToSchemaGroup 转换为群组对象
func (a Group) ToSchemaGroup() *schema.Group {
	item := new(schema.Group)
	structure.Copy(a, item)
	return item
}

// Groups 群组实体列表
type Groups []*Group

// ToSchemaGroups 转换为群组对象列表
func (a Groups) ToSchemaGroups() []*schema.Group {
	list := make([]*schema.Group, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaGroup()
	}

	return list
}
