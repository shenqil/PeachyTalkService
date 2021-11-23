package service

import (
	"context"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"github.com/google/wire"
	"strings"
)

// FriendSet 注入好友
var FriendSet = wire.NewSet(wire.Struct(new(Friend), "*"))

// Friend 好友管理
type Friend struct {
	UserModel       *repo.User
	UserFriendModel *repo.UserFriend
}

// Query 查询数据
func (a *Friend) Query(ctx context.Context, params schema.UserFriendQueryParam, opts ...schema.UserFriendQueryOptions) (*schema.FriendListQueryResult, error) {
	result, err := a.UserFriendModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	}

	userInfoResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserIDs: result.Data.ToFriendIDs(params.UserID),
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

	if oldUserFriend.ID == "" {
		err = a.UserFriendModel.Create(ctx, userFriend)
	} else {
		err = a.UserFriendModel.Update(ctx, userFriend.ID, userFriend)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	//	合并数据
	if userFriend.Status1 == schema.FriendNone {
		userFriend.Status1 = oldUserFriend.Status1
	}

	if userFriend.Status2 == schema.FriendNone {
		userFriend.Status2 = oldUserFriend.Status2
	}

	return &userFriend, nil
}

// Delete 删除好友
func (a *Friend) Delete(ctx context.Context, fromUserID, toUserID string) (*schema.UserFriend, error) {
	// 得到表结构数据
	userFriend := a.getUserFriendInfo(fromUserID, toUserID, schema.FriendUnsubscribe, schema.FriendUnsubscribe)

	err := a.UserFriendModel.Update(ctx, userFriend.ID, userFriend)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &userFriend, nil
}
