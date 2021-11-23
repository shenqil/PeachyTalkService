package mqttTopic

import (
	"ginAdmin/internal/app/mqttApi"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
)

var _ ITopic = (*Topic)(nil)

// TopicSet 注入Topic
var TopicSet = wire.NewSet(wire.Struct(new(Topic), "*"), wire.Bind(new(ITopic), new(*Topic)))

// ITopic 注册主题
type ITopic interface {
	Register(cli mqtt.Client) error
}

// Topic 主题管理器
type Topic struct {
	UserAPI     *mqttApi.User
	ManifestAPI *mqttApi.Manifest
	FriendAPI   *mqttApi.Friend
}

func (a *Topic) Register(cli mqtt.Client) error {
	a.RegisterAPI(cli)
	return nil
}
