package schema

import "time"

// Role 角色对象
type Role struct {
	ID          string      `json:"id"`                                    // 唯一标识
	Name        string      `json:"name" binding:"required"`               // 角色名称
	Sequence    int         `json:"sequence"`                              // 排序值
	Memo        string      `json:"memo"`                                  // 备注
	Status      int         `json:"status" binding:"required,max=2,min=1"` // 状态(1:启用 2:禁用)
	Creator     string      `json:"creator"`                               // 创建者
	CreatedAt   time.Time   `json:"createdAt"`                             // 创建时间
	UpdatedAt   time.Time   `json:"updatedAt"`                             // 更新时间
	RoleRouters RoleRouters `json:"roleRouters"`                           // 路由资源
}

// RoleQueryParam 查询条件
type RoleQueryParam struct {
	PaginationParam
	IDs        []string `form:"-"`          // 唯一标识列表
	ExcludeIDs []string `form:"-"`          // 排除的唯一标识列表
	Name       string   `form:"-"`          // 角色名称
	QueryValue string   `form:"queryValue"` // 模糊查询
	UserID     string   `form:"-"`          // 用户ID
	Status     int      `form:"status"`     // 状态(1:启用 2:禁用)
}

// RoleQueryOptions 查询可选参数项
type RoleQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// RoleQueryResult 查询结果
type RoleQueryResult struct {
	Data       Roles             `json:"list"`
	PageResult *PaginationResult `json:"pagination"`
}

// Roles 角色对象列表
type Roles []*Role

// ToNames 获取角色名称列表
func (a Roles) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}
	return names
}

// ToMap 转换为键值存储
func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// ToIDs 转换为id
func (a Roles) ToIDs() []string {
	ids := make([]string, len(a))
	for i, item := range a {
		ids[i] = item.ID
	}
	return ids
}

// --------------------------------- RoleShow -------------------------------

// RoleShow 角色显示项
type RoleShow struct {
	ID              string            `json:"id"`                                    // 唯一标识
	Name            string            `json:"name" binding:"required"`               // 角色名称
	Sequence        int               `json:"sequence"`                              // 排序值
	Memo            string            `json:"memo"`                                  // 备注
	Status          int               `json:"status" binding:"required,max=2,min=1"` // 状态(1:启用 2:禁用)
	Creator         string            `json:"creator"`                               // 创建者
	CreatedAt       time.Time         `json:"createdAt"`                             // 创建时间
	RouterResources []*RouterResource `json:"routerResources"`                       // 角色管理的路由资源
}

// RoleShows 角色显示列表
type RoleShows []*RoleShow

// RoleShowQueryResult 用户显示项查询结果
type RoleShowQueryResult struct {
	Data       RoleShows
	PageResult *PaginationResult
}
