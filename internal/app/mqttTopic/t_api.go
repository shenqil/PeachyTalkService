package mqttTopic

import (
	"fmt"
	"ginAdmin/internal/app/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

var prefix string

// RegisterAPI 注册 api
func (a *Topic) RegisterAPI(cli mqtt.Client) {
	prefix = config.C.MQTT.TopicPrefix
	if token := cli.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	subscribe(cli, "manifest/get/#", 0, a.ManifestAPI.Get)
	subscribe(cli, "user/get/#", 0, a.UserAPI.Get)
	subscribe(cli, "user/getToken/#", 0, a.UserAPI.GetToken)
}

// 订阅一个主题
func subscribe(c mqtt.Client, topic string, qos byte, handle mqtt.MessageHandler) {
	if token := c.Subscribe(fmt.Sprintf("%s/%s", prefix, topic), qos, handle); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
