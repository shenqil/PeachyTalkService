package mqttApi

import (
	"context"
	"encoding/json"
	"fmt"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"ginAdmin/pkg/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"time"
)

// FriendSet 注入好友
var FriendSet = wire.NewSet(wire.Struct(new(Friend), "*"))

// Friend 好友关系
type Friend struct {
	FriendSrc *service.Friend
	UserSrv   *service.User
}

// Search 查询需要添加的好友
func (a *Friend) Search(client mqtt.Client, msg mqtt.Message) {
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
	keywords := string(msg.Payload())

	// 查询到数据
	result, err := a.FriendSrc.Search(ctx, keywords)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, result)
}

// MyFriendList 查询我的好友
func (a *Friend) MyFriendList(client mqtt.Client, msg mqtt.Message) {
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
	userID := string(msg.Payload())

	// 拿到好友列表
	result, err := a.FriendSrc.MyFriendList(ctx, userID)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, result)
}

// QuasiFriendList 获取准好友列表
func (a *Friend) QuasiFriendList(client mqtt.Client, msg mqtt.Message) {
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
	userID := string(msg.Payload())

	// 拿到好友列表
	result, err := a.FriendSrc.QuasiFriendList(ctx, userID)
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

	quasiFriend := &schema.QuasiFriend{
		Info:   friendInfo.ToFriendInfo(),
		Status: userFriend,
	}

	// 推送消息给双方
	err = friendsChange(client, userName, quasiFriend)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}
	err = friendsChange(client, friendInfo.UserName, quasiFriend)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
}

// Ignore 忽略好友
func (a *Friend) Ignore(client mqtt.Client, msg mqtt.Message) {
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
	userFriend, err := a.FriendSrc.Ignore(ctx, params.FormUserID, params.ToUserID)
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

	quasiFriend := &schema.QuasiFriend{
		Info:   friendInfo.ToFriendInfo(),
		Status: userFriend,
	}

	// 推送消息给双方
	err = friendsChange(client, userName, quasiFriend)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}
	err = friendsChange(client, friendInfo.UserName, quasiFriend)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
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
	userFriend, err := a.FriendSrc.Delete(ctx, params.FormUserID, params.ToUserID)
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

	quasiFriend := &schema.QuasiFriend{
		Info:   friendInfo.ToFriendInfo(),
		Status: userFriend,
	}

	// 推送消息给双方
	err = friendsChange(client, userName, quasiFriend)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}
	err = friendsChange(client, friendInfo.UserName, quasiFriend)
	if err != nil {
		replyError(client, userName, msgID, err.Error())
		return
	}

	replySuccess(client, userName, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
}
