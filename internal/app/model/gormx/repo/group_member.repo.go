package repo

import (
	"context"
	"ginAdmin/internal/app/model/gormx/entity"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// GroupMemberSet 注入群组成员GroupMember
var GroupMemberSet = wire.NewSet(wire.Struct(new(GroupMember), "*"))

// GroupMember 群组成员关联
type GroupMember struct {
	DB *gorm.DB
}

func (a *GroupMember) getQueryOption(opts ...schema.GroupMemberQueryOptions) schema.GroupMemberQueryOptions {
	var opt schema.GroupMemberQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	return opt
}

// Query 查询数据
func (a *GroupMember) Query(ctx context.Context, params schema.GroupMemberQueryParam, opts ...schema.GroupMemberQueryOptions) (*schema.GroupMemberQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetGroupMemberDB(ctx, a.DB)

	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.GroupMembers
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := &schema.GroupMemberQueryResult{
		Data:       list.ToSchemaGroupMembers(),
		PageResult: pr,
	}

	return result, nil
}

// Get 查询指定数据
func (a *GroupMember) Get(ctx context.Context, id string, opts ...schema.GroupQueryOptions) (*schema.GroupMember, error) {
	db := entity.GetGroupMemberDB(ctx, a.DB).Where("id=?", id)

	var item entity.GroupMember
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaGroupMember(), nil
}

// Create 创建数据
func (a *GroupMember) Create(ctx context.Context, item schema.GroupMember) error {
	sitem := entity.SchemaGroupMember(item)

	result := entity.GetGroupMemberDB(ctx, a.DB).Create(sitem.ToGroupMember())

	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *GroupMember) Update(ctx context.Context, id string, item schema.GroupMember) error {
	sitem := entity.SchemaGroupMember(item)
	result := entity.GetGroupMemberDB(ctx, a.DB).Where("id=?", id).Updates(sitem.ToGroupMember())
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *GroupMember) Delete(ctx context.Context, id string) error {
	result := entity.GetGroupMemberDB(ctx, a.DB).Where("id=?", id).Delete(&entity.GroupMember{})

	return errors.WithStack(result.Error)
}
