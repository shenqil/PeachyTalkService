package mqttApi

import (
	"errors"
	"fmt"
	"ginAdmin/pkg/util/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

// 回复成功消息
func replySuccess(client mqtt.Client, userName string, msgId string, payload interface{}) {
	buf, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	topic := fmt.Sprintf("%s/reply/success/%s", userName, msgId)
	client.Publish(topic, 0, false, buf)
}

// 回复错误消息
func replyError(client mqtt.Client, userName string, msgId string, err string) {
	topic := fmt.Sprintf("%s/reply/error/%s", userName, msgId)
	client.Publish(topic, 0, false, err)
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
