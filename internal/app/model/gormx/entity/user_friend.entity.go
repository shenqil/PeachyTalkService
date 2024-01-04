package entity

import (
	"PeachyTalkService/internal/app/schema"
	"PeachyTalkService/pkg/util/structure"
	"context"
	"time"

	"gorm.io/gorm"
)

// GetUserFriendDB 获取用户好友关系储存
func GetUserFriendDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(UserFriend))
}

// SchemaUserFriend 用户好友
type SchemaUserFriend schema.UserFriend

// ToUserFriend 转换为好友实体
func (a SchemaUserFriend) ToUserFriend() *UserFriend {
	item := new(UserFriend)
	structure.Copy(a, item)
	return item
}

// UserFriend 用户好友关联实体
type UserFriend struct {
	ID        string         `gorm:"column:id;primaryKey;size:72;"`
	UserID1   string         `gorm:"column:user_id1;size:36;index;default:'';not null;"` // 用户1内码
	UserID2   string         `gorm:"column:user_id2;size:36;index;default:'';not null;"` // 用户2内码
	Status1   int            `gorm:"column:status1;size:36;index;default:0;not null;"`   // 用户1好友状态
	Status2   int            `gorm:"column:status2;size:36;index;default:0;not null;"`   // 用户2好友状态
	CreatedAt time.Time      `gorm:"column:created_at;index;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;"`
}

// ToSchemaUserFriend 转换为用户好友对象
func (a UserFriend) ToSchemaUserFriend() *schema.UserFriend {
	item := new(schema.UserFriend)
	structure.Copy(a, item)
	return item
}

// UserFriends 用户好友关联列表
type UserFriends []*UserFriend

// ToSchemaUserFriends 转换为用户好友对象列表
func (a UserFriends) ToSchemaUserFriends() []*schema.UserFriend {
	list := make([]*schema.UserFriend, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserFriend()
	}
	return list
}
