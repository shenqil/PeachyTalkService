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
		return err
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

// friendsChange 好友变动
func friendsChange(client mqtt.Client, userName string, payload interface{}) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("IMClient/%s/friend/change/%d", userName, time.Now().UnixMilli())
	token := client.Publish(topic, 0, false, buf)
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
