package service

import (
	"context"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"ginAdmin/pkg/util/uuid"
	"github.com/google/wire"
)

// RouterResourceSet 注入RouterResource
var RouterResourceSet = wire.NewSet(wire.Struct(new(RouterResource), "*"))

// RouterResource 路由资源
type RouterResource struct {
	RouterResourceModel *repo.RouterResource
}

// Query 查询数据
func (a *RouterResource) Query(ctx context.Context, params schema.RouterResourceQueryParam, opts ...schema.RouterResourceQueryOptions) (*schema.RouterResourceResult, error) {
	return a.RouterResourceModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *RouterResource) Get(ctx context.Context, id string, opts ...schema.RouterResourceQueryOptions) (*schema.RouterResource, error) {
	return a.RouterResourceModel.Get(ctx, id, opts...)
}

func (a *RouterResource) checkName(ctx context.Context, name string) error {
	result, err := a.RouterResourceModel.Query(ctx, schema.RouterResourceQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Name:            name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("名称已存在")
	}

	return nil
}

// Create 创建数据
func (a *RouterResource) Create(ctx context.Context, item schema.RouterResource) (*schema.IDResult, error) {
	err := a.checkName(ctx, item.Name)
	if err != nil {
		return nil, err
	}

	item.ID = uuid.MustString()
	err = a.RouterResourceModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

// Update 更新数据
func (a *RouterResource) Update(ctx context.Context, id string, item schema.RouterResource) error {
	oldItem, err := a.RouterResourceModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		if err := a.checkName(ctx, item.Name); err != nil {
			return err
		}
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	return a.RouterResourceModel.Update(ctx, id, item)
}

// Delete 删除数据
func (a *RouterResource) Delete(ctx context.Context, id string) error {
	oldItem, err := a.RouterResourceModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.RouterResourceModel.Delete(ctx, id)
}

// UpdateStatus 更新状态
func (a *RouterResource) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.RouterResourceModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.RouterResourceModel.UpdateStatus(ctx, id, status)
}
