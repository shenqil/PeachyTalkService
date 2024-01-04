//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package app

import (
	"ginAdmin/internal/app/api"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/router"
	"ginAdmin/internal/app/service"
	"github.com/google/wire"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		// mock.MockSet,
		InitGormDB,
		repo.RepoSet,
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
