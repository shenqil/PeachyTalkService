package repo

import (
	"context"
	"ginAdmin/internal/app/model/gormx/entity"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// RouterResource 注入RouterResource
var RouterResourceSet = wire.NewSet(wire.Struct(new(RouterResource), "*"))

// RouterResource 路由资源
type RouterResource struct {
	DB *gorm.DB
}

func (a *RouterResource) getQueryOption(opts ...schema.RouterResourceQueryOptions) schema.RouterResourceQueryOptions {
	var opt schema.RouterResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *RouterResource) Query(ctx context.Context, params schema.RouterResourceQueryParam, opts ...schema.RouterResourceQueryOptions) (*schema.RouterResourceResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetRouterResourceDB(ctx, a.DB)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR memo LIKE ? OR path LIKE ?", v, v, v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status = ?", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("created_at", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.RouterResources
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, err
	}

	qr := &schema.RouterResourceResult{
		PageResult: pr,
		Data:       list.ToSchemaRouterResources(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *RouterResource) Get(ctx context.Context, id string, opts ...schema.RouterResourceQueryOptions) (*schema.RouterResource, error) {
	var item entity.RouterResource

	ok, err := FindOne(ctx, entity.GetRouterResourceDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item.ToSchemaRouterResource(), nil
}

// Create 创建数据
func (a *RouterResource) Create(ctx context.Context, item schema.RouterResource) error {
	eitem := entity.SchemaRouterResource(item).ToRouterResource()
	result := entity.GetRouterResourceDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *RouterResource) Update(ctx context.Context, id string, item schema.RouterResource) error {
	eitem := entity.SchemaRouterResource(item).ToRouterResource()
	result := entity.GetRouterResourceDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *RouterResource) Delete(ctx context.Context, id string) error {
	result := entity.GetRouterResourceDB(ctx, a.DB).Where("id=?", id).Delete(entity.RouterResource{})
	return errors.WithStack(result.Error)
}

// UpdateStatus 更新状态
func (a *RouterResource) UpdateStatus(ctx context.Context, id string, status int) error {
	result := entity.GetRouterResourceDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
