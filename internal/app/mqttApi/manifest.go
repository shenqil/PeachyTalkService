package mqttApi

import (
	"context"
	"fmt"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"ginAdmin/pkg/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
)

// ManifestSet 注入Manifest
var ManifestSet = wire.NewSet(wire.Struct(new(Manifest), "*"))

// Manifest 主清单
type Manifest struct {
	UserSrv *service.User
}

// Get 查询指定用户的主清单
func (a Manifest) Get(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户名称和消息id
	userName, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	manifest := new(schema.Manifest)

	// 拿到用户信息
	result, err := a.UserSrv.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{},
		UserName:        userName,
		Status:          1,
	})
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}
	ids := result.Data.ToIDs()
	if len(ids) != 1 {
		logger.WithContext(ctx).Fatalf(fmt.Sprintf("取到用户id数量：%d != 1", len(ids)))
		return
	}

	manifest.UserID = ids[0]

	replySuccess(client, userName, msgID, manifest)
}
