package mqttApi

import (
	"github.com/google/wire"
)

// MQTTApiSet 注入API
var MQTTApiSet = wire.NewSet(
	UserSet,
	ManifestSet,
	FriendSet,
)
