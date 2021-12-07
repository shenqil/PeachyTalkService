package schema

import (
	"time"
)

// -------------------------------- Group -------------------------------

// Group 群组对象
type Group struct {
	ID        string    `json:"id"`        // ID
	GroupName string    `json:"groupName"` // 群名称
	Brief     string    `json:"brief"`     // 简介
	Avatar    string    `json:"avatar"`    // 头像
	Owner     string    `json:"owner"`     // 拥有者
	Creator   string    `json:"creator"`   // 创建者
	CreatedAt time.Time `json:"createdAt"`
	MemberIDs []string  `json:"memberIDs"` // 成员
}

// GroupQueryParam 查询条件
type GroupQueryParam struct {
	PaginationParam
	Owner  string   `form:"owner"` // 拥有者
	IDs    []string `form:"-"`     // 群组id 列表
	UserID string   `json:"-"`     // 用户ID
}

// GroupQueryOptions 查询可选参数项
type GroupQueryOptions struct {
	OrderFields []*OrderField
}

// GroupQueryResult 查询结果
type GroupQueryResult struct {
	Data       Groups
	PageResult *PaginationResult
}

// Groups 群组列表
type Groups []*Group

// ToGroupIDs 转为为群组id 列表
func (a Groups) ToGroupIDs() []string {
	list := make([]string, len(a))
	for i, group := range a {
		list[i] = group.ID
	}

	return list
}

// --------------------------------- GroupMember--------------------------

// GroupMember 群组成员关联对象
type GroupMember struct {
	ID      string `json:"id"`      // ID
	GroupID string `json:"groupId"` // 群组ID
	UserID  string `json:"userId"`  // 用户ID
}

// GroupMemberQueryParam 查询条件
type GroupMemberQueryParam struct {
	PaginationParam
	UserID   string   `form:"userId"` // 用户id
	GroupIDs []string `form:"-"`      // 群组id
}

// GroupMemberQueryOptions 查询可选参数项
type GroupMemberQueryOptions struct {
	OrderFields []*OrderField
}

// GroupMemberQueryResult 查询结果
type GroupMemberQueryResult struct {
	Data       GroupMembers
	PageResult *PaginationResult
}

// GroupMemberChangesInfo 成员变动信息
type GroupMemberChangesInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GroupMemberChangesParam 成员变动参数
type GroupMemberChangesParam struct {
	FromID   string `json:"fromId"`   // 操作人id
	FromName string `json:"fromName"` // 操作人名称
	GroupID  string `json:"groupId"`  // 群ID
	List     []*GroupMemberChangesInfo
}

// GroupMembers 群组成员关联列表对象
type GroupMembers []*GroupMember

// ToMemberIDs 转换为成员ID列表
func (a GroupMembers) ToMemberIDs() []string {
	list := make([]string, len(a))
	for i, member := range a {
		list[i] = member.UserID
	}

	return list
}

// ToGroupIDMap 转换为群组id映射
func (a GroupMembers) ToGroupIDMap() map[string]GroupMembers {
	m := make(map[string]GroupMembers)
	for _, member := range a {
		m[member.GroupID] = append(m[member.GroupID], member)
	}
	return m
}
