package mqttApi

import (
	"context"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"ginAdmin/pkg/logger"
	"ginAdmin/pkg/util/json"
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
func (a *User) Get(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 拿到用户信息
	item, err := a.UserSrv.Get(ctx, userID)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	replySuccess(client, userID, msgID, item.CleanSecure())
}

// GetToken 获取Token
func (a *User) GetToken(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 生成一个Token
	tokenInfo, err := a.LoginSrv.GenerateToken(ctx, userID)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	replySuccess(client, userID, msgID, tokenInfo)
}

// GetUserInfo  根据 id 查找用户信息
func (a *User) GetUserInfo(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 解析参数
	var ids []string
	err = json.Unmarshal(msg.Payload(), &ids)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	result, err := a.UserSrv.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{
			Pagination: false,
		},
		Status:  1,
		UserIDs: ids,
	})

	replySuccess(client, userID, msgID, result.Data.ToFriendInfo())
}
