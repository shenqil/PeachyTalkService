package adapter

import (
	"context"
	"fmt"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/schema"
	"ginAdmin/pkg/logger"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/google/wire"
)

var _ persist.Adapter = (*CasbinAdapter)(nil)

// CasbinAdapterSet 注入 CasbinAdapter
var CasbinAdapterSet = wire.NewSet(wire.Struct(new(CasbinAdapter), "*"), wire.Bind(new(persist.Adapter), new(*CasbinAdapter)))

// CasbinAdapter casbin适配器
type CasbinAdapter struct {
	RoleModel       *repo.Role
	RouterModel     *repo.RouterResource
	RoleRouterModel *repo.RoleRouter
	UserModel       *repo.User
	UserRoleModel   *repo.UserRole
}

// LoadPolicy 从存储加载所有策略规则
func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	ctx := context.Background()
	err := a.LoadRolePolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}

	err = a.LoadUserPolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}

	return nil
}

// 加载角色策略(p,role_id,path,method)
func (a *CasbinAdapter) LoadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	// 1.查询所有角色
	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	} else if len(roleResult.Data) == 0 {
		return nil
	}

	// 2.查询角色与路由资源对应表
	roleRouterResult, err := a.RoleRouterModel.Query(ctx, schema.RoleRouterQueryParam{})
	if err != nil {
		return err
	}
	mRoleRouters := roleRouterResult.Data.ToRoleIDMap()

	// 3.查询出所有路由资源
	routerResourceResult, err := a.RouterModel.Query(ctx, schema.RouterResourceQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	}
	mRouterResource := routerResourceResult.Data.ToMap()

	// 4.加载角色策略
	for _, item := range roleResult.Data {
		mCahce := make(map[string]struct{})
		if roleRoters, ok := mRoleRouters[item.ID]; ok {
			for _, routerID := range roleRoters.ToRouterIDs() {
				if rr, ok := mRouterResource[routerID]; ok {
					if rr.Path == "" || rr.Method == "" {
						continue
					} else if _, ok := mCahce[rr.Path+rr.Method]; ok {
						continue
					}
					mCahce[rr.Path+rr.Method] = struct{}{}
					line := fmt.Sprintf("p,%s,%s,%s", item.ID, rr.Path, rr.Method)
					persist.LoadPolicyLine(line, m)
				}
			}
		}
	}

	return nil
}

// 加载用户策略(g,user_id,role_id)
func (a *CasbinAdapter) LoadUserPolicy(ctx context.Context, m casbinModel.Model) error {
	userResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	} else if len(userResult.Data) > 0 {
		userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{})
		if err != nil {
			return err
		}

		mUserRoles := userRoleResult.Data.ToUserIDMap()
		for _, uitem := range userResult.Data {
			if urs, ok := mUserRoles[uitem.ID]; ok {
				for _, ur := range urs {
					line := fmt.Sprintf("g,%s,%s", ur.UserID, ur.RoleID)
					persist.LoadPolicyLine(line, m)
				}
			}
		}
	}
	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
