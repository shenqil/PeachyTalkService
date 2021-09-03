package service

import (
	"context"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/errors"
	"ginAdmin/pkg/util/uuid"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
)

// RoleSet 注入Role
var RoleSet = wire.NewSet(wire.Struct(new(Role), "*"))

// Role 角色管理
type Role struct {
	Enforcer            *casbin.SyncedEnforcer
	TransModel          *repo.Trans
	RoleModel           *repo.Role
	UserModel           *repo.User
	RoleRouterModel     *repo.RoleRouter
	RouterResourceModel *repo.RouterResource
}

// Query 查询数据
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	return a.RoleModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	item, err := a.RoleModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	roleRouter, err := a.RoleRouterModel.Query(ctx, schema.RoleRouterQueryParam{
		RoleID: id,
	})
	if err != nil {
		return nil, err
	}

	item.RoleRouters = roleRouter.Data

	return item, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) (*schema.IDResult, error) {
	err := a.checkName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.ID = uuid.MustString()
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		return a.RoleModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (a *Role) checkName(ctx context.Context, item schema.Role) error {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		Name:            item.Name,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("角色名称已经存在")
	}
	return nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, id string, item schema.Role) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		err := a.checkName(ctx, item)
		if err != nil {
			return err
		}
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		return a.RoleModel.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, id string) error {
	oldItem, err := a.RoleModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	userResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		RoleIDs:         []string{id},
	})
	if err != nil {
		return err
	} else if userResult.PageResult.Total > 0 {
		return errors.New400Response("该角色已被赋予用户，不允许删除")
	}

	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		return a.RoleModel.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.RoleModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.RoleModel.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}
