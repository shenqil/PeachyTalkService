package repo

import "github.com/google/wire"

// RepoSet model 注入
var RepoSet = wire.NewSet(
	DemoSet,
	MenuSet,
	RoleSet,
	TransSet,
	UserRoleSet,
	UserSet,
	RoleRouterSet,
	RouterResourceSet,
)
