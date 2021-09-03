package schema

import "time"

// 路由资源
type RouterResource struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`               // 资源名称
	Memo      string    `json:"memo"`                                  // 资源备注
	Method    string    `json:"method" binding:"required"`             // 资源请求方式(支持正则)
	Path      string    `json:"path" binding:"required"`               // 资源请求路径（支持/:id匹配）
	Status    int       `json:"status" binding:"required,max=2,min=1"` // 状态(1:启用 2:禁用)
	Creator   string    `json:"creator"`                               // 创建者
	CreatedAt time.Time `json:"createdAt"`                             // 创建时间
	UpdatedAt time.Time `json:"updatedAt"`
}

// RouterResourceQueryParam 查询条件
type RouterResourceQueryParam struct {
	PaginationParam
	IDs        []string `form:"-"`          // 唯一标识列表
	Name       string   `form:"name"`       // 名称
	RoleID     string   `form:"roleId"`     // 角色ID
	QueryValue string   `form:"queryValue"` // 模糊查询
	Status     int      `form:"status"`     // 启用状态
}

// RouterResourceQueryOptions 查询可选参数
type RouterResourceQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// RouterResourceResult 路由资源对象查询结果
type RouterResourceResult struct {
	Data       []*RouterResource
	PageResult *PaginationResult
}
