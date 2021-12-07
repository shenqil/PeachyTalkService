package repo

import (
	"context"
	"ginAdmin/internal/app/model/gormx/entity"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// GroupSet 注入Group
var GroupSet = wire.NewSet(wire.Struct(new(Group), "*"))

// Group 群组储存
type Group struct {
	DB *gorm.DB
}

func (a *Group) getQueryOption(opts ...schema.GroupQueryOptions) schema.GroupQueryOptions {
	var opt schema.GroupQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	return opt
}

// Query 查询数据
func (a *Group) Query(ctx context.Context, params schema.GroupQueryParam, opts ...schema.GroupQueryOptions) (*schema.GroupQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetGroupDB(ctx, a.DB)
	if v := params.Owner; v != "" {
		db = db.Where("owner=?", v)
	}

	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := params.UserID; v != "" {
		subQuery := entity.GetGroupMemberDB(ctx, a.DB).
			Select("group_id").
			Where("user_id=?", v)
		db = db.Where("id IN (?)", subQuery)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("created_at", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.Groups
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.GroupQueryResult{
		Data:       list.ToSchemaGroups(),
		PageResult: pr,
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Group) Get(ctx context.Context, id string, opts ...schema.GroupQueryOptions) (*schema.Group, error) {
	var item entity.Group

	ok, err := FindOne(ctx, entity.GetGroupDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaGroup(), nil
}

// Create 创建数据
func (a *Group) Create(ctx context.Context, item schema.Group) error {
	sitem := entity.SchemaGroup(item)

	result := entity.GetGroupDB(ctx, a.DB).Create(sitem.ToGroup())

	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *Group) Update(ctx context.Context, id string, item schema.Group) error {
	sitem := entity.SchemaGroup(item)

	result := entity.GetGroupDB(ctx, a.DB).Updates(sitem.ToGroup())

	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *Group) Delete(ctx context.Context, id string) error {
	result := entity.GetGroupDB(ctx, a.DB).Where("id=?", id).Delete(&entity.Group{})
	return errors.WithStack(result.Error)
}
