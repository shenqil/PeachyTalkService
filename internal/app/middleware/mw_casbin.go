package middleware

import (
	"ginAdmin/internal/app/config"
	"ginAdmin/internal/app/ginx"
	"ginAdmin/pkg/errors"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.C.Casbin
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		if b, err := enforcer.Enforce(ginx.GetUserID(c), p, m); err != nil {
			ginx.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			ginx.ResError(c, errors.ErrNoPerm)
			return
		}
		c.Next()
	}
}
