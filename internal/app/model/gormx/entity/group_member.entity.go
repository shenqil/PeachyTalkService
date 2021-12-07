package entity

import (
	"context"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/util/structure"
	"gorm.io/gorm"
)

// GetGroupMemberDB 获取群组成员关联数据
func GetGroupMemberDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(GroupMember))
}

// SchemaGroupMember 群组成员关联对象
type SchemaGroupMember schema.GroupMember

// ToGroupMember 转换为群组成员关联实体
func (a SchemaGroupMember) ToGroupMember() *GroupMember {
	item := new(GroupMember)
	structure.Copy(a, item)
	return item
}

// GroupMember 群组成员关联实体
type GroupMember struct {
	ID      string `gorm:"column:id;primaryKey;size:36;"`           // ID
	GroupID string `gorm:"column:group_id;size:36;not null;index;"` // 群组ID
	UserID  string `gorm:"column:user_id;size:36;not null;index;"`  // 用户ID
}

// ToSchemaGroupMember 转换为群组成员关联对象
func (a *GroupMember) ToSchemaGroupMember() *schema.GroupMember {
	item := new(schema.GroupMember)
	structure.Copy(a, item)
	return item
}

// GroupMembers 群组成员关联列表实体
type GroupMembers []*GroupMember

// ToSchemaGroupMembers 转换为群组成员列表关联实体
func (a GroupMembers) ToSchemaGroupMembers() schema.GroupMembers {
	list := make(schema.GroupMembers, len(a))
	for i, member := range a {
		list[i] = member.ToSchemaGroupMember()
	}

	return list
}
