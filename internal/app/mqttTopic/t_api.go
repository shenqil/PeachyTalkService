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
	// 公用
	subscribe(cli, "manifest/get/#", 0, a.ManifestAPI.Get)
	subscribe(cli, "user/get/#", 0, a.UserAPI.Get)
	subscribe(cli, "user/getToken/#", 0, a.UserAPI.GetToken)

	// 好友
	subscribe(cli, "friend/search/#", 0, a.FriendAPI.Search)
	subscribe(cli, "friend/myFriends/#", 0, a.FriendAPI.MyFriendList)
	subscribe(cli, "friend/quasiFriends/#", 0, a.FriendAPI.QuasiFriendList)
	subscribe(cli, "friend/add/#", 0, a.FriendAPI.Add)
	subscribe(cli, "friend/ignore/#", 0, a.FriendAPI.Ignore)
	subscribe(cli, "friend/delete/#", 0, a.FriendAPI.Delete)

	//	群组
	subscribe(cli, "group/query/#", 0, a.GroupAPI.Query)
	subscribe(cli, "group/create/#", 0, a.GroupAPI.Create)
	subscribe(cli, "group/delete/#", 0, a.GroupAPI.Delete)
	subscribe(cli, "group/update/#", 0, a.GroupAPI.Update)
	subscribe(cli, "group/addMembers/#", 0, a.GroupAPI.AddMembers)
	subscribe(cli, "group/delMembers/#", 0, a.GroupAPI.DelMembers)
	subscribe(cli, "group/exitGroup/#", 0, a.GroupAPI.ExitGroup)
}

// 订阅一个主题
func subscribe(c mqtt.Client, topic string, qos byte, handle mqtt.MessageHandler) {
	if token := c.Subscribe(fmt.Sprintf("%s/%s", prefix, topic), qos, handle); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
