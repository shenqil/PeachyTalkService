package service

import (
	"PeachyTalkService/internal/app/model/gormx/repo"
	"PeachyTalkService/internal/app/schema"
	"PeachyTalkService/pkg/errors"
	"context"
	"strings"
	"time"

	"github.com/google/wire"
)

// FriendSet 注入好友
var FriendSet = wire.NewSet(wire.Struct(new(Friend), "*"))

// Friend 好友管理
type Friend struct {
	UserModel       *repo.User
	UserFriendModel *repo.UserFriend
}

// Search 查询指定数据
func (a *Friend) Search(ctx context.Context, keywords string, opts ...schema.UserFriendQueryOptions) (*schema.FriendInfo, error) {
	userInfoResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		PreciseSearch:   keywords,
		Status:          1,
		PaginationParam: schema.PaginationParam{Pagination: false},
	})
	if err != nil {
		return nil, err
	}

	if len(userInfoResult.Data) == 0 {
		return nil, nil
	}

	friendInfo := userInfoResult.Data.ToFriendInfo()

	return friendInfo[0], nil
}

// MyFriendList 获取我的好友
func (a *Friend) MyFriendList(ctx context.Context, userID string, opts ...schema.UserFriendQueryOptions) (*schema.FriendListQueryResult, error) {
	result, err := a.UserFriendModel.MyFriendList(ctx, userID)
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	} else if len(result.Data) == 0 {
		return nil, nil
	}

	userInfoResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserIDs:         result.Data.ToFriendIDs(userID),
		PaginationParam: schema.PaginationParam{Pagination: false},
	})
	if err != nil {
		return nil, err
	}

	UserFriendResult := &schema.FriendListQueryResult{
		Data:       userInfoResult.Data.ToFriendInfo(),
		PageResult: result.PageResult,
	}

	return UserFriendResult, nil
}

// QuasiFriendList 获取准好友列表
func (a *Friend) QuasiFriendList(ctx context.Context, userID string, opts ...schema.UserFriendQueryOptions) (*schema.QuasiFriendQueryResult, error) {
	result, err := a.UserFriendModel.QuasiFriendList(ctx, userID, opts...)
	if err != nil {
		return nil, err
	} else if result == nil || len(result.Data) == 0 {
		return nil, nil
	}

	userInfoResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserIDs:         result.Data.ToFriendIDs(userID),
		PaginationParam: schema.PaginationParam{Pagination: false},
	})
	if err != nil {
		return nil, err
	}

	friendList := userInfoResult.Data.ToFriendInfo()

	// 合并数据
	quasiFriends := make(schema.QuasiFriends, len(friendList))
	for i, item := range friendList {
		var status *schema.UserFriend

		// 查找好友状态
		for _, item2 := range result.Data {
			if item2.UserID1 == item.ID || item2.UserID2 == item.ID {
				status = item2
				break
			}
		}
		quasiFriends[i] = &schema.QuasiFriend{
			Info:   item,
			Status: status,
		}
	}

	return &schema.QuasiFriendQueryResult{
		Data:       quasiFriends,
		PageResult: nil,
	}, nil
}

func (a *Friend) getUserFriendInfo(id1, id2 string, status1, status2 int) schema.UserFriend {
	userFriend := schema.UserFriend{}
	if strings.Compare(id1, id2) == 1 {
		userFriend.UserID1 = id1
		userFriend.UserID2 = id2
		userFriend.ID = userFriend.UserID1 + userFriend.UserID2
		userFriend.Status1 = status1
		userFriend.Status2 = status2
	} else {
		userFriend.UserID1 = id2
		userFriend.UserID2 = id1
		userFriend.ID = userFriend.UserID1 + userFriend.UserID2
		userFriend.Status1 = status2
		userFriend.Status2 = status1
	}

	return userFriend
}

// Add 添加好友
func (a *Friend) Add(ctx context.Context, fromUserID, toUserID string) (*schema.UserFriend, error) {
	// 得到表结构数据
	userFriend := a.getUserFriendInfo(fromUserID, toUserID, schema.FriendSubscribe, schema.FriendNone)

	oldUserFriend, err := a.UserFriendModel.Get(ctx, userFriend.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if oldUserFriend == nil {
		// 创建
		err = a.UserFriendModel.Create(ctx, userFriend)
	} else {
		// 更新

		// 如果对方忽略，重新添加时置为FriendUnsubscribe
		if fromUserID == oldUserFriend.UserID1 {
			if oldUserFriend.Status2 == schema.FriendIgnore {
				userFriend.Status2 = schema.FriendUnsubscribe
			}
		} else {
			if oldUserFriend.Status1 == schema.FriendIgnore {
				userFriend.Status1 = schema.FriendUnsubscribe
			}
		}

		err = a.UserFriendModel.Update(ctx, userFriend.ID, userFriend)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result, err := a.UserFriendModel.Get(ctx, userFriend.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

// Ignore 忽略对方本地添加
func (a *Friend) Ignore(ctx context.Context, fromUserID, toUserID string) (*schema.UserFriend, error) {
	userFriend := a.getUserFriendInfo(fromUserID, toUserID, schema.FriendIgnore, schema.FriendNone)

	err := a.UserFriendModel.Update(ctx, userFriend.ID, userFriend)
	if err != nil {
		return nil, err
	}

	result, err := a.UserFriendModel.Get(ctx, userFriend.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

// Delete 删除好友
func (a *Friend) Delete(ctx context.Context, fromUserID, toUserID string) (*schema.UserFriend, error) {
	// 得到表结构数据
	userFriend := a.getUserFriendInfo(fromUserID, toUserID, schema.FriendUnsubscribe, schema.FriendUnsubscribe)

	err := a.UserFriendModel.Delete(ctx, userFriend.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	userFriend.UpdatedAt = time.Now()

	return &userFriend, nil
}
