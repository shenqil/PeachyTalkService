package app

import (
	"fmt"
	"ginAdmin/internal/app/config"
	"ginAdmin/internal/app/mqttTopic"
	"ginAdmin/pkg/util/hash"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("DefaultPublish:TOPIC: %s - MSG: %s", msg.Topic(), msg.Payload())
}

// InitMQTTSever  初始化 mqtt服务
func InitMQTTSever(t mqttTopic.Topic) func() {
	cfg := config.C.MQTT
	root := config.C.Root

	addr := fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port)
	opts := mqtt.NewClientOptions().AddBroker(addr).SetClientID(cfg.ClientID)

	opts.SetUsername(root.UserName)
	opts.SetPassword(hash.MD5String(root.Password))
	opts.SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second)
	opts.SetPingTimeout(time.Duration(cfg.PingTimeout) * time.Second)
	// 设置消息回调处理函数
	opts.SetDefaultPublishHandler(f)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 注册所有主题
	t.Register(c)

	return func() {
		// 断开连接
		c.Disconnect(250)
	}
}
