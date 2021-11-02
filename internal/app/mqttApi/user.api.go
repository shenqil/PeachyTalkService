package mqttApi

import (
	"context"
	"ginAdmin/internal/app/service"
	"ginAdmin/pkg/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

// User 用户管理
type User struct {
	UserSrv  *service.User
	LoginSrv *service.Login
}

// Get 查询指定数据
func (a User) Get(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userName, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 拿到用户信息
	item, err := a.UserSrv.GetByUserName(ctx, userName)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, item.CleanSecure())
}

// GetToken 获取Token
func (a User) GetToken(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userName, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 拿到用户信息
	item, err := a.UserSrv.GetByUserName(ctx, userName)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	// 生成一个Token
	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, item.ID)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, tokenInfo)
}
