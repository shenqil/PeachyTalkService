package bucket

import "github.com/google/wire"

// BucketSet model 注入
var BucketSet = wire.NewSet(
	AvatarSet,
	//UserFile,
)
