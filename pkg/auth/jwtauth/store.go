package jwtauth

import (
	"context"
	"time"
)

// Storer 令牌储存接口
type Storer interface {
	//	储存令牌数据，并指定到期时间
	Set(ctx context.Context, tokenString string, expiration time.Duration) error
	//	检查令牌是否存在
	Check(ctx context.Context, tokenString string) (bool, error)
	//	关闭储存
	Close() error
}
