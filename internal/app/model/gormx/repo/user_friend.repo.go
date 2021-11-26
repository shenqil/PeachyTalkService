package repo

import (
	"context"
	"ginAdmin/internal/app/model/gormx/entity"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// UserFriendSet 注入UserFriend
var UserFriendSet = wire.NewSet(wire.Struct(new(UserFriend), "*"))

// UserFriend 用户好友储存
type UserFriend struct {
	DB *gorm.DB
}

func (a *UserFriend) getQueryOption(opts ...schema.UserFriendQueryOptions) schema.UserFriendQueryOptions {
	var opt schema.UserFriendQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *UserFriend) Query(ctx context.Context, userID string, opts ...schema.UserFriendQueryOptions) (*schema.UserFriendQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetUserFriendDB(ctx, a.DB)
	if userID != "" {
		db = db.Where("user_id1=? OR user_id2=?", userID, userID)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.UserFriends
	pr, err := WrapPageQuery(ctx, db, schema.PaginationParam{Pagination: false}, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserFriendQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserFriends(),
	}
	return qr, nil
}

// MyFriendList 获取我的好友列表
func (a *UserFriend) MyFriendList(ctx context.Context, userID string) (*schema.UserFriendQueryResult, error) {
	db := entity.GetUserFriendDB(ctx, a.DB)
	db = db.Where("user_id1=? OR user_id2=?", userID, userID).Where("status1=? AND status2=?", schema.FriendSubscribe, schema.FriendSubscribe)

	var list entity.UserFriends
	pr, err := WrapPageQuery(ctx, db, schema.PaginationParam{Pagination: false}, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserFriendQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserFriends(),
	}
	return qr, nil
}

// QuasiFriendList 获取准好友列表
func (a *UserFriend) QuasiFriendList(ctx context.Context, userID string, opts ...schema.UserFriendQueryOptions) (*schema.UserFriendQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetUserFriendDB(ctx, a.DB)
	db = db.Where("user_id1=? OR user_id2=?", userID, userID).Where("status1<>? or status2<>?", schema.FriendSubscribe, schema.FriendSubscribe)

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("updated_at", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.UserFriends
	pr, err := WrapPageQuery(ctx, db, schema.PaginationParam{Pagination: false}, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserFriendQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserFriends(),
	}
	return qr, nil
}

// Get 查询指定数据
func (a *UserFriend) Get(ctx context.Context, id string, opts ...schema.UserRoleQueryOptions) (*schema.UserFriend, error) {
	db := entity.GetUserFriendDB(ctx, a.DB).Where("id=?", id)
	var item entity.UserFriend
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUserFriend(), nil
}

// Create 创建数据
func (a *UserFriend) Create(ctx context.Context, item schema.UserFriend) error {
	eItem := entity.SchemaUserFriend(item).ToUserFriend()
	result := entity.GetUserFriendDB(ctx, a.DB).Create(eItem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *UserFriend) Update(ctx context.Context, id string, item schema.UserFriend) error {
	eItem := entity.SchemaUserFriend(item).ToUserFriend()
	result := entity.GetUserFriendDB(ctx, a.DB).Where("id=?", id).Updates(eItem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *UserFriend) Delete(ctx context.Context, id string) error {
	result := entity.GetUserFriendDB(ctx, a.DB).Where("id=?", id).Delete(entity.UserFriend{})
	return errors.WithStack(result.Error)
}
