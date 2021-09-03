package schema

// RoleRouter 角色资源
type RoleRouter struct {
	ID       string `json:"id"`       // 唯一标识
	RoleID   string `json:"roleId"`   // 角色ID
	RouterID string `json:"routerId"` // 路由资源
}

// RoleRouterQueryParam 查询条件
type RoleRouterQueryParam struct {
	PaginationParam
	RoleID  string   // 角色ID
	RoleIDs []string // 角色ID列表
}

// RoleRouterQueryOptions 查询可选参数项
type RoleRouterQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// RoleRouterQueryResult 查询结果
type RoleRouterQueryResult struct {
	Data       RoleRouters
	PageResult *PaginationResult
}

// RoleRouter 角色路由列表
type RoleRouters []*RoleRouter

// ToMap 转换为map
func (a RoleRouters) ToMap() map[string]*RoleRouter {
	m := make(map[string]*RoleRouter)
	for _, item := range a {
		m[item.RouterID] = item
	}
	return m
}

// ToRouterIDs 转换为路由资源ID列表
func (a RoleRouters) ToRouterIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RouterID
	}
	return list
}

// ToRoleIDMap 转换为角色ID映射
func (a RoleRouters) ToRoleIDMap() map[string]RoleRouters {
	m := make(map[string]RoleRouters)
	for _, item := range a {
		m[item.RoleID] = append(m[item.RoleID], item)
	}
	return m
}
