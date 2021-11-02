//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package app

import (
	"ginAdmin/internal/app/api"
	"ginAdmin/internal/app/model/gormx/repo"
	"ginAdmin/internal/app/module/adapter"
	"ginAdmin/internal/app/mqttApi"
	"ginAdmin/internal/app/mqttTopic"
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
		InitCasbin,
		InitGinEngine,
		service.ServiceSet,
		api.APISet,
		router.RouterSet,
		adapter.CasbinAdapterSet,
		mqttTopic.TopicSet,
		mqttApi.MQTTApiSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
