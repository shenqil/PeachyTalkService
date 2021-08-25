// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package app

import (
	"ginAdmin/internal/app/api"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/module/adapter"
	"ginAdmin/internal/app/router"
	"ginAdmin/internal/app/service"
)

// Injectors from wire.go:

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	auther, cleanup, err := InitAuth()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := InitGormDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	role := &repo.Role{
		DB: db,
	}
	user := &repo.User{
		DB: db,
	}
	userRole := &repo.UserRole{
		DB: db,
	}
	casbinAdapter := &adapter.CasbinAdapter{
		RoleModel:     role,
		UserModel:     user,
		UserRoleModel: userRole,
	}
	syncedEnforcer, cleanup3, err := InitCasbin(casbinAdapter)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	demo := &repo.Demo{
		DB: db,
	}
	serviceDemo := &service.Demo{
		DemoModel: demo,
	}
	apiDemo := &api.Demo{
		DemoSrv: serviceDemo,
	}
	menu := &repo.Menu{
		DB: db,
	}
	login := &service.Login{
		Auth:          auther,
		UserModel:     user,
		UserRoleModel: userRole,
		RoleModel:     role,
		MenuModel:     menu,
	}
	apiLogin := &api.Login{
		LoginSrv: login,
	}
	trans := &repo.Trans{
		DB: db,
	}
	serviceMenu := &service.Menu{
		TransModel: trans,
		MenuModel:  menu,
	}
	apiMenu := &api.Menu{
		MenuSrv: serviceMenu,
	}
	serviceRole := &service.Role{
		Enforcer:   syncedEnforcer,
		TransModel: trans,
		RoleModel:  role,
		UserModel:  user,
	}
	apiRole := &api.Role{
		RoleSrv: serviceRole,
	}
	serviceUser := &service.User{
		Enforcer:      syncedEnforcer,
		TransModel:    trans,
		UserModel:     user,
		UserRoleModel: userRole,
		RoleModel:     role,
	}
	apiUser := &api.User{
		UserSrv: serviceUser,
	}
	routerRouter := &router.Router{
		Auth:           auther,
		CasbinEnforcer: syncedEnforcer,
		DemoAPI:        apiDemo,
		LoginAPI:       apiLogin,
		MenuAPI:        apiMenu,
		RoleAPI:        apiRole,
		UserAPI:        apiUser,
	}
	engine := InitGinEngine(routerRouter)
	injector := &Injector{
		Engine:         engine,
		Auth:           auther,
		CasbinEnforcer: syncedEnforcer,
		MenuBll:        serviceMenu,
	}
	return injector, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
