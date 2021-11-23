package schema

// 枚举好友状态
const (
	FriendSubscribe   = 1
	FriendUnsubscribe = 2
	FriendRefuse      = 3
	FriendNone        = 0
)

// -------------------------- FriendInfo------------------------------

// FriendInfo 好友信息
type FriendInfo struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	RealName string `json:"real_name"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// FriendListQueryResult 查询结果
type FriendListQueryResult struct {
	Data       FriendList
	PageResult *PaginationResult
}

// FriendList 用户好友列表
type FriendList []*FriendInfo

//  ------------------------- UserFriend-------------------------------

// UserFriend 用户好友
type UserFriend struct {
	ID      string `json:"id"`
	UserID1 string `json:"userID1"` // 用户内码1
	UserID2 string `json:"userID1"` // 用户内码2
	Status1 int    `json:"status1"` // 用户1好友状态
	Status2 int    `json:"status2"` // 用户2好友状态
}

// UserFriendQueryParam 查询条件
type UserFriendQueryParam struct {
	UserID string `form:"userID" json:"userID"` // 用户ID
	PaginationParam
}

// UserFriendOperateParam 好友操作参数
type UserFriendOperateParam struct {
	FormUserID string `json:"formUserId"`
	ToUserID   string `json:"toUserId"`
}

// UserFriendQueryOptions 查询可选参数项
type UserFriendQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserFriendQueryResult 查询结果
type UserFriendQueryResult struct {
	Data       UserFriends
	PageResult *PaginationResult
}

// UserFriends 用户好友列表
type UserFriends []*UserFriend

// ToFriendIDs 转换为好友id列表
func (a UserFriends) ToFriendIDs(userID string) []string {
	list := make([]string, len(a))
	for i, item := range a {
		if item.UserID1 == userID {
			list[i] = item.UserID2
		} else {
			list[i] = item.UserID1
		}
	}
	return list
}
