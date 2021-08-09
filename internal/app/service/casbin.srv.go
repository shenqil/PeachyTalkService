package service

import (
	"context"
	"ginAdmin/internal/app/config"
	"ginAdmin/pkg/logger"
	"github.com/casbin/casbin/v2"
)

var chCasbinPolicy chan *chCasbinPolicyItem

type chCasbinPolicyItem struct {
	ctx context.Context
	e   *casbin.SyncedEnforcer
}

func init() {
	chCasbinPolicy = make(chan *chCasbinPolicyItem, 1)
	go func() {
		for item := range chCasbinPolicy {
			err := item.e.LoadPolicy()
			if err != nil {
				logger.WithContext(item.ctx).Errorf("The Load casbin policy error: %s", err.Error())
			}
		}
	}()
}

// loadCasbinPolicy 异步加载casbin权限策略
func LoadCasbinPolicy(ctx context.Context, e *casbin.SyncedEnforcer) {
	if !config.C.Casbin.Enable {
		return
	}

	if len(chCasbinPolicy) > 0 {
		logger.WithContext(ctx).Infof("The load casbin policy is already in the wait quue")
	}

	chCasbinPolicy <- &chCasbinPolicyItem{
		ctx: ctx,
		e:   e,
	}
}
