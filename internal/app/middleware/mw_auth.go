package middleware

import (
	"PeachyTalkService/internal/app/config"
	"PeachyTalkService/internal/app/contextx"
	"PeachyTalkService/internal/app/ginx"
	"PeachyTalkService/pkg/auth"
	"PeachyTalkService/pkg/errors"
	"PeachyTalkService/pkg/logger"

	"github.com/gin-gonic/gin"
)

// 包装用户身份验证上下文
func wrapUserAuthContext(c *gin.Context, userID string) {
	ginx.SetUserID(c, userID)
	ctx := contextx.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, config.C.Root.UserName)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(c.Request.Context(), ginx.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config.C.IsDebugMode() {
					wrapUserAuthContext(c, config.C.Root.UserName)
					c.Next()
					return
				}
			}
			ginx.ResError(c, errors.ErrInvalidToken)
			return
		}
		wrapUserAuthContext(c, userID)
		c.Next()
	}
}
