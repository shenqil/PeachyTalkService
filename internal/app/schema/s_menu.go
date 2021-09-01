package schema

import (
	"ginAdmin/pkg/util/json"
	"strings"
	"time"
)

// Menu 菜单对象
type Menu struct {
	ID         string    `json:"id"`                                        // 唯一标识
	Name       string    `json:"name" binding:"required"`                   // 菜单名称
	Sequence   int       `json:"sequence"`                                  // 排序值
	Icon       string    `json:"icon"`                                      // 菜单图标
	Router     string    `json:"router"`                                    // 访问路由
	ParentID   string    `json:"parentId"`                                  // 父级ID
	ParentPath string    `json:"parentPath"`                                // 父级路径
	ShowStatus int       `json:"showStatus" binding:"required,max=2,min=1"` // 显示状态(1:显示 2:隐藏)
	Status     int       `json:"status" binding:"required,max=2,min=1"`     // 状态(1:启用 2:禁用)
	Memo       string    `json:"memo"`                                      // 备注
	Creator    string    `json:"creator"`                                   // 创建者
	CreatedAt  time.Time `json:"createdAt"`                                 // 创建时间
	UpdatedAt  time.Time `json:"updatedAt"`                                 // 更新时间
}

func (a *Menu) String() string {
	return json.MarshalToString(a)
}

// MenuQueryParam 查询条件
type MenuQueryParam struct {
	PaginationParam
	IDs              []string `form:"-"`          // 唯一标识列表
	Name             string   `form:"-"`          // 菜单名称
	PrefixParentPath string   `form:"-"`          // 父级路径(前缀模糊查询)
	QueryValue       string   `form:"queryValue"` // 模糊查询
	ParentID         *string  `form:"parentID"`   // 父级内码
	ShowStatus       int      `form:"showStatus"` // 显示状态(1:显示 2:隐藏)
	Status           int      `form:"status"`     // 状态(1:启用 2:禁用)
}

// MenuQueryOptions 查询可选参数项
type MenuQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// MenuQueryResult 查询结果
type MenuQueryResult struct {
	Data       Menus
	PageResult *PaginationResult
}

// Menus 菜单列表
type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	return a[i].Sequence > a[j].Sequence
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// ToMap 转换为键值映射
func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// SplitParentIDs 拆分父级路径的唯一标识列表
func (a Menus) SplitParentIDs() []string {
	idList := make([]string, 0, len(a))
	mIDList := make(map[string]struct{})

	for _, item := range a {
		if _, ok := mIDList[item.ID]; ok || item.ParentPath == "" {
			continue
		}

		for _, pp := range strings.Split(item.ParentPath, "/") {
			if _, ok := mIDList[pp]; ok {
				continue
			}
			idList = append(idList, pp)
			mIDList[pp] = struct{}{}
		}
	}

	return idList
}

// ToTree 转换为菜单树
func (a Menus) ToTree() MenuTrees {
	list := make(MenuTrees, len(a))
	for i, item := range a {
		list[i] = &MenuTree{
			ID:         item.ID,
			Name:       item.Name,
			Icon:       item.Icon,
			Router:     item.Router,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			Sequence:   item.Sequence,
			ShowStatus: item.ShowStatus,
			Status:     item.Status,
		}
	}
	return list.ToTree()
}

// ----------------------------------------MenuTree--------------------------------------

// MenuTree 菜单树
type MenuTree struct {
	ID         string     `yaml:"-" json:"id"`                                  // 唯一标识
	Name       string     `yaml:"name" json:"name"`                             // 菜单名称
	Icon       string     `yaml:"icon" json:"icon"`                             // 菜单图标
	Router     string     `yaml:"router,omitempty" json:"router"`               // 访问路由
	ParentID   string     `yaml:"-" json:"parentId"`                            // 父级ID
	ParentPath string     `yaml:"-" json:"parentPath"`                          // 父级路径
	Sequence   int        `yaml:"sequence" json:"sequence"`                     // 排序值
	ShowStatus int        `yaml:"-" json:"showStatus"`                          // 显示状态(1:显示 2:隐藏)
	Status     int        `yaml:"-" json:"status"`                              // 状态(1:启用 2:禁用)
	Children   *MenuTrees `yaml:"children,omitempty" json:"children,omitempty"` // 子级树
}

// MenuTrees 菜单树列表
type MenuTrees []*MenuTree

// ToTree 转换为树形结构
func (a MenuTrees) ToTree() MenuTrees {
	mi := make(map[string]*MenuTree)
	for _, item := range a {
		mi[item.ID] = item
	}

	var list MenuTrees
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				children := MenuTrees{item}
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		} else {
			list = append(list, item)
		}
	}
	return list
}
