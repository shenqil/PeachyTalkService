package mqttApi

import (
	"context"
	"fmt"
	"ginAdmin/internal/app/schema"
	"ginAdmin/internal/app/service"
	"ginAdmin/pkg/errors"
	"ginAdmin/pkg/logger"
	"ginAdmin/pkg/util/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"time"
)

// GroupSet 注入群组
var GroupSet = wire.NewSet(wire.Struct(new(Group), "*"))

// Group 群组
type Group struct {
	GroupSrc *service.Group
	UserSrv  *service.User
}

// Query 查询群组
func (a *Group) Query(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	// 解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	result, err := a.GroupSrc.Query(ctx, schema.GroupQueryParam{
		UserID: userID,
	})

	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	replySuccess(client, userID, msgID, result)
}

// Create 创建群组
func (a *Group) Create(client mqtt.Client, msg mqtt.Message) {
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
	var params schema.Group
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	// 创建成功
	result, err := a.GroupSrc.Create(ctx, params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
	}

	// 拿到最新数据
	groupInfo, err := a.GroupSrc.Get(ctx, result.ID)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
	}

	// 推送给所有相关的人
	for _, id := range groupInfo.MemberIDs {
		err := groupChange(client, id, groupInfo.ID, "create", groupInfo)
		if err != nil {
			logger.WithContext(ctx).Fatalf(err.Error())
		}
	}

	// 创建成功
	replySuccess(client, userID, msgID, result.ID)
}

// Delete 删除群组 群主专用
func (a *Group) Delete(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	//	解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 解析参数
	groupID := string(msg.Payload())
	if groupID == "" {
		replyError(client, userID, msgID, errors.ErrNotFound.Error())
	}

	// 删除群组
	groupInfo, err := a.GroupSrc.Delete(ctx, groupID, userID)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
	}

	// 推送给所有相关的人
	for _, id := range groupInfo.MemberIDs {
		err := groupChange(client, id, groupID, "delete", groupInfo)
		if err != nil {
			logger.WithContext(ctx).Fatalf(err.Error())
		}
	}

	// 删除成功
	replySuccess(client, userID, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
}

// Update 更新数据 群主专用
func (a *Group) Update(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	//	解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 解析参数
	var params schema.Group
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	//	更新数据
	groupInfo, err := a.GroupSrc.Update(ctx, params.ID, userID, params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
	}

	// 推送给所有相关的人
	for _, id := range groupInfo.MemberIDs {
		err := groupChange(client, id, groupInfo.ID, "update", groupInfo)
		if err != nil {
			logger.WithContext(ctx).Fatalf(err.Error())
		}
	}

	// 操作成功
	replySuccess(client, userID, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
}

// AddMembers 添加成员
func (a *Group) AddMembers(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	//	解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 解析参数
	var params schema.GroupMemberChangesParam
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	//	操作
	result, memberIDs, err := a.GroupSrc.AddMembers(ctx, params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	// 更新被操作的人员
	params.List = result

	// 通知所有成员
	if result != nil && memberIDs != nil {
		for _, id := range memberIDs {
			err := groupChange(client, id, params.GroupID, "addMembers", params)
			if err != nil {
				logger.WithContext(ctx).Fatalf(err.Error())
			}
		}
	}

	replySuccess(client, userID, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
}

// DelMembers 删除成员
func (a *Group) DelMembers(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	//	解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 解析参数
	var params schema.GroupMemberChangesParam
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	//	操作
	result, memberIDs, err := a.GroupSrc.DelMembers(ctx, params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}
	// 更新被操作的人员
	params.List = result

	// 通知所有成员
	if result != nil && memberIDs != nil {
		for _, id := range memberIDs {
			err := groupChange(client, id, params.GroupID, "delMembers", params)
			if err != nil {
				logger.WithContext(ctx).Fatalf(err.Error())
			}
		}
	}

	replySuccess(client, userID, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
}

// ExitGroup 退出群聊
func (a *Group) ExitGroup(client mqtt.Client, msg mqtt.Message) {
	// 创建一个上下文
	ctx := logger.NewTraceIDContext(context.Background(), msg.Topic())
	ctx = logger.NewTagContext(ctx, "__MQTT__")

	//	解析用户id和消息id
	userID, msgID, err := parseUserNameAndMsgIDWithTopic(msg.Topic())
	if err != nil {
		logger.WithContext(ctx).Fatalf(err.Error())
		return
	}

	// 解析参数
	var params schema.GroupMemberChangesParam
	err = json.Unmarshal(msg.Payload(), &params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	//	退出群聊
	memberIDs, err := a.GroupSrc.ExitGroup(ctx, params)
	if err != nil {
		replyError(client, userID, msgID, err.Error())
		return
	}

	// 通知所有成员
	if memberIDs != nil {
		for _, id := range memberIDs {
			err := groupChange(client, id, params.GroupID, "exitGroup", params)
			if err != nil {
				logger.WithContext(ctx).Fatalf(err.Error())
			}
		}
	}

	replySuccess(client, userID, msgID, fmt.Sprintf("%d", time.Now().UnixMilli()))
}
