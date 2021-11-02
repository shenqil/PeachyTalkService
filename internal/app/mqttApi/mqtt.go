package mqttApi

import (
	"errors"
	"fmt"
	"ginAdmin/internal/app/config"
	"ginAdmin/pkg/util/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

// 回复消息
func replyJSON(client mqtt.Client, userId string, msgId string, qos byte, retained bool, payload interface{}) {
	buf, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	topic := fmt.Sprintf("%s/reply/%s/%s", config.C.MQTT.TopicPrefix, userId, msgId)
	client.Publish(topic, qos, retained, buf)
}

// 回复错误
func replyErr(client mqtt.Client, userId string, msgId string, err string) {
	topic := fmt.Sprintf("%s/replyError/%s/%s", config.C.MQTT.TopicPrefix, userId, msgId)
	client.Publish(topic, 0, false, err)
}

// 解析 topic 里面的用户id消息id
func parseUserIDAndMsgIDWithTopic(topic string) (userID, MsgID string, err error) {
	ary := strings.Split(topic, "/")

	l := len(ary)
	if len(ary) < 4 {
		return "", "", errors.New("错误的topic:topic分级小于4")
	}

	MsgID = ary[l-1]
	userID = ary[l-2]

	return userID, MsgID, nil
}
