//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package app

import (
	"PeachyTalkService/internal/app/api"
	"PeachyTalkService/internal/app/model/gormx/repo"
	"PeachyTalkService/internal/app/model/miniox/bucket"
	"PeachyTalkService/internal/app/router"
	"PeachyTalkService/internal/app/service"
	"github.com/google/wire"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		// mock.MockSet,
		InitGormDB,
		InitMinio,
		repo.RepoSet,
		bucket.BucketSet,
		InitAuth,
		InitGinEngine,
		service.ServiceSet,
		api.APISet,
		router.RouterSet,
		//mqttTopic.TopicSet,
		//mqttApi.MQTTApiSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
