package schema

// StatusText 定义状态文本
type StatusText string

func (t StatusText) String() string {
	return string(t)
}

// 定义HTTP状态文本常量
const (
	OKStatus    StatusText = "OK"
	ErrorStatus StatusText = "ERROR"
	FailStatus  StatusText = "FAIL"
)

// StatusResult 响应状态
type StatusResult struct {
	Status StatusText `json:"status"` // 状态(OK)
}

// ErrorResult 响应错误
type ErrorResult struct {
	Error ErrorItem `json:"error"` // 错误项
}

// ErrorItem 响应错误项
type ErrorItem struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

// ListResult 响应列表数据
type ListResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

// PaginationResult 分页查询结果
type PaginationResult struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

// PaginationParam 分页查询条件
type PaginationParam struct {
	Pagination bool `form:"-"`                                                     // 是否使用分页查询
	OnlyCount  bool `form:"-"`                                                     // 是否仅查询count
	Current    int  `form:"current,default=1" json:"current"`                      // 当前页
	PageSize   int  `form:"pageSize,default=10" binding:"max=100" json:"pageSize"` // 页大小
}

// GetCurrent 获取当前页
func (a PaginationParam) GetCurrent() int {
	return a.Current
}

// GetPageSize 获取页大小
func (a PaginationParam) GetPageSize() int {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 100
	}
	return pageSize
}

// OrderDirection 排序方向
type OrderDirection int

const (
	// OrderByASC 升序排序
	OrderByASC OrderDirection = 1
	// OrderByDESC 降序排序
	OrderByDESC OrderDirection = 2
)

// NewOrderFieldWithKeys 创建排序字段(默认升序排序)，可指定不同key的排序规则
// keys 需要排序的key
// directions 排序规则，按照key的索引指定，索引默认从0开始
func NewOrderFieldWithKeys(keys []string, directions ...map[int]OrderDirection) []*OrderField {
	m := make(map[int]OrderDirection)
	if len(directions) > 0 {
		m = directions[0]
	}

	fields := make([]*OrderField, len(keys))
	for i, key := range keys {
		d := OrderByASC
		if v, ok := m[i]; ok {
			d = v
		}

		fields[i] = NewOrderField(key, d)
	}

	return fields
}

// NewOrderFields 创建排序字段列表
func NewOrderFields(orderFields ...*OrderField) []*OrderField {
	return orderFields
}

// NewOrderField 创建排序字段
func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

// OrderField 排序字段
type OrderField struct {
	Key       string         // 字段名(字段名约束为小写蛇形)
	Direction OrderDirection // 排序方向
}

// NewIDResult 创建响应唯一标识实例
func NewIDResult(id string) *IDResult {
	return &IDResult{
		ID: id,
	}
}

// IDsRequest 请求的id列表
type IDsRequest struct {
	IDs []string `json:"ids"`
}

// IDResult 响应唯一标识
type IDResult struct {
	ID string `json:"id"`
}
