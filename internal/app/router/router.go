package router

import (
	"ginAdmin/internal/app/api"
	"ginAdmin/pkg/auth"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var _ IRouter = (*Router)(nil)

// RouterSet 注入router
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

// IRouter 注册路由
type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router 路由管理器
type Router struct {
	Auth           auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
	DemoAPI        *api.Demo
	LoginAPI       *api.Login
	MenuAPI        *api.Menu
	RoleAPI        *api.Role
	UserAPI        *api.User
	RouterResource *api.RouterResource
	IMApi          *api.IM
}

func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	a.RegisterIMAPI(app)
	return nil
}

func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}
