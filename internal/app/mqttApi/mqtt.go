package mqttApi

import (
	"errors"
	"fmt"
	"ginAdmin/pkg/util/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"time"
)

// 通用回复成功消息
func replySuccess(client mqtt.Client, userName string, msgId string, payload interface{}) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	topic := fmt.Sprintf("IMClient/%s/reply/success/%s", userName, msgId)
	token := client.Publish(topic, 0, false, buf)
	token.Wait()

	return token.Error()
}

// 通用回复错误消息
func replyError(client mqtt.Client, userName string, msgId string, err string) error {
	topic := fmt.Sprintf("IMClient/%s/reply/error/%s", userName, msgId)
	token := client.Publish(topic, 0, false, err)
	token.Wait()

	return token.Error()
}

// 发送添加好友消息
func joinFriend(client mqtt.Client, formUserName string, toUserName string) error {
	topic := fmt.Sprintf("IMClient/%s/friend/join/%s", toUserName, formUserName)
	token := client.Publish(topic, 1, true, fmt.Sprintf("%d", time.Now().UnixMilli()))
	token.Wait()

	return token.Error()
}

// 好友添加成功
func joinFriendSuccess(client mqtt.Client, formUserName string, toUserName string) error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())

	topic := fmt.Sprintf("IMClient/%s/friend/joinSuccess/%s", toUserName, formUserName)
	token := client.Publish(topic, 1, true, timestamp)
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}

	topic = fmt.Sprintf("IMClient/%s/friend/joinSuccess/%s", formUserName, toUserName)
	token = client.Publish(topic, 1, true, timestamp)
	token.Wait()

	return token.Error()
}

// 删除好友
func deleteFriend(client mqtt.Client, formUserName string, toUserName string) error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())

	topic := fmt.Sprintf("IMClient/%s/friend/deleteFriend/%s", toUserName, formUserName)
	token := client.Publish(topic, 1, true, timestamp)
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}

	topic = fmt.Sprintf("IMClient/%s/friend/deleteFriend/%s", formUserName, toUserName)
	token = client.Publish(topic, 1, true, timestamp)
	token.Wait()

	return token.Error()
}

// 解析 topic 里面的用户名称和消息id
func parseUserNameAndMsgIDWithTopic(topic string) (userName, MsgID string, err error) {
	ary := strings.Split(topic, "/")

	l := len(ary)
	if len(ary) < 4 {
		return "", "", errors.New("错误的topic:topic分级小于4")
	}

	MsgID = ary[l-1]
	userName = ary[l-2]

	return userName, MsgID, nil
}
