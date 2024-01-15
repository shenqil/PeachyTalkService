package schema

import (
	"PeachyTalkService/internal/app/config"
	"PeachyTalkService/pkg/util/hash"
	"PeachyTalkService/pkg/util/json"
	"context"
	"time"
)

// GetRootUser 获取root用户
func GetRootUser() *User {
	user := config.C.Root
	return &User{
		ID:       user.UserName,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: hash.MD5String(user.Password),
	}
}

// CheckIsRootUser 检查是否是root用户
func CheckIsRootUser(ctx context.Context, userID string) bool {
	return GetRootUser().ID == userID
}

// User 用户对象
type User struct {
	ID          string    `json:"id"`                                    // 唯一标识
	UserName    string    `json:"userName" binding:"required"`           // 用户名
	RealName    string    `json:"realName" binding:"required"`           // 真实姓名
	Password    string    `json:"password"`                              // 密码
	Avatar      string    `json:"avatar"`                                // 头像
	Gender      int       `json:"gender"`                                // 性别(1:男 0:女)
	DateOfBirth string    `json:"dateOfBirth"`                           // 出生日期
	Phone       string    `json:"phone"`                                 // 手机号
	Email       string    `json:"email"`                                 // 邮箱
	Status      int       `json:"status" binding:"required,max=2,min=1"` // 用户状态(1:启用 2:停用)
	Creator     string    `json:"creator"`                               // 创建者
	CreatedAt   time.Time `json:"createdAt"`                             // 创建时间
}

func (a *User) String() string {
	return json.MarshalToString(a)
}

// CleanSecure 清理安全数据
func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

// ToFriendInfo 转换为好友信息
func (a *User) ToFriendInfo() *FriendInfo {
	return &FriendInfo{
		ID:       a.ID,
		UserName: a.UserName,
		RealName: a.RealName,
		Avatar:   a.Avatar,
		Phone:    a.Phone,
		Email:    a.Email,
	}
}

// UserQueryParam 查询条件
type UserQueryParam struct {
	PaginationParam
	UserName      string   `form:"userName"`   // 用户名
	QueryValue    string   `form:"queryValue"` // 模糊查询
	Status        int      `form:"status"`     // 用户状态(1:启用 2:停用)
	PreciseSearch string   `form:"-"`          // 精确查询查询
	UserIDs       []string `form:"-"`          // 用户ID列表
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserQueryResult 查询结果
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// Users 用户对象列表
type Users []*User

// ToIDs 转换为唯一标识列表
func (a Users) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

// CleanSecure 清理安全数据
func (a Users) CleanSecure() Users {
	for _, user := range a {
		user.CleanSecure()
	}
	return a
}

// ToFriendInfo 转换为好友信息
func (a Users) ToFriendInfo() FriendList {
	list := make(FriendList, len(a))
	for i, item := range a {
		list[i] = item.ToFriendInfo()
	}

	return list
}
