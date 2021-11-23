package mqttApi

import (
	"context"
	"encoding/json"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"ginAdmin/pkg/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
)

// FriendSet 注入好友
var FriendSet = wire.NewSet(wire.Struct(new(Friend), "*"))

// Friend 好友关系
type Friend struct {
	FriendSrc *service.Friend
	UserSrv   *service.User
}

// Query 查询我的好友
func (a *Friend) Query(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userName, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	//	解析参数
	var params schema.UserFriendQueryParam
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	// 拿到好友列表
	result, err := a.FriendSrc.Query(ctx, params)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, result)
}

// Add 添加好友
func (a *Friend) Add(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userName, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	//	解析参数
	var params schema.UserFriendOperateParam
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	//	执行操作
	userFriend, err := a.FriendSrc.Add(ctx, params.FormUserID, params.ToUserID)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	//	拿到需要添加的好友信息
	friendInfo, err := a.UserSrv.Get(ctx, params.ToUserID)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	// 需要对方同意
	if userFriend.Status1 != schema.FriendSubscribe || userFriend.Status2 != schema.FriendSubscribe {
		err := joinFriend(client, friendInfo.UserName, userName)
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 双方已同意
	err = joinFriendSuccess(client, userName, friendInfo.UserName)
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
	}
}

// Delete 删除好友
func (a *Friend) Delete(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userName, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	//	解析参数
	var params schema.UserFriendOperateParam
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	//	执行操作
	_, err = a.FriendSrc.Delete(ctx, params.FormUserID, params.ToUserID)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	//	拿到需要添加的好友信息
	friendInfo, err := a.UserSrv.Get(ctx, params.ToUserID)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	err = deleteFriend(client, userName, friendInfo.UserName)
}
