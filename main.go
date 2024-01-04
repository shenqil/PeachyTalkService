package main

import (
	"PeachyTalkService/internal/app"
	"PeachyTalkService/pkg/logger"
	"context"
)

// VERSION 版本号
var VERSION = "1.0.0"

func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTagContext(context.Background(), "__main__")
	app.Run(ctx,
		app.SetConfigFile("./configs/config.toml"),
		app.SetVersion(VERSION),
	)
}
