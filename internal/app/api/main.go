package api

import "github.com/google/wire"

// APISet 注入API
var APISet = wire.NewSet(
	DemoSet,
	LoginSet,
	MenuSet,
	RoleSet,
	UserSet,
)
