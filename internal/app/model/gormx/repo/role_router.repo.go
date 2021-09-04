package repo

import (
	"context"
	"ginAdmin/internal/app/model/gormx/entity"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var RoleRouterSet = wire.NewSet(wire.Struct(new(RoleRouter), "*"))

// RoleRouter 角色路由储存
type RoleRouter struct {
	DB *gorm.DB
}

func (a *RoleRouter) getQueryOption(opts ...schema.RoleRouterQueryOptions) schema.RoleRouterQueryOptions {
	var opt schema.RoleRouterQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *RoleRouter) Query(ctx context.Context, params schema.RoleRouterQueryParam, opts ...schema.RoleRouterQueryOptions) (*schema.RoleRouterQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetRoleRouterDB(ctx, a.DB)
	if v := params.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}
	if v := params.RouterID; v != "" {
		db = db.Where("router_id=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		db = db.Where("role_id in (?)", v)
	}

	if opt.OrderFields != nil {
		db = db.Order(ParseOrder(opt.OrderFields))
	}

	var list entity.RoleRouters
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleRouterQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoleRouters(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *RoleRouter) Get(ctx context.Context, id string, opts ...schema.RoleRouterQueryOptions) (*schema.RoleRouter, error) {
	var item entity.RoleRouter
	db := entity.GetRoleRouterDB(ctx, a.DB).Where("id=?", id)
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRoleRouter(), nil
}

// Create 创建数据
func (a *RoleRouter) Create(ctx context.Context, item schema.RoleRouter) error {
	eitem := entity.SchemaRoleRouter(item).ToRoleRouter()
	result := entity.GetRoleRouterDB(ctx, a.DB).Where("id=?", item.ID).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *RoleRouter) Update(ctx context.Context, id string, item schema.RoleRouter) error {
	eitem := entity.SchemaRoleRouter(item).ToRoleRouter()
	result := entity.GetRoleRouterDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *RoleRouter) Delete(ctx context.Context, id string) error {
	result := entity.GetRoleRouterDB(ctx, a.DB).Where("id=?", id).Delete(schema.RoleRouter{})
	return errors.WithStack(result.Error)
}

// DeleteByRoleID 根据角色ID 删除数据
func (a *RoleRouter) DeleteByRoleUD(ctx context.Context, roleID string) error {
	result := entity.GetRoleRouterDB(ctx, a.DB).Where("role_id=?", roleID).Delete(schema.RoleRouter{})
	return errors.WithStack(result.Error)
}
