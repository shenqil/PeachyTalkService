package service

import (
	"PeachyTalkService/internal/app/model/miniox/bucket"
	"PeachyTalkService/internal/app/schema"
	"context"
	"github.com/google/wire"
)

// FileSet 注入File
var FileSet = wire.NewSet(wire.Struct(new(File), "*"))

// File 文件
type File struct {
	AvatarModel *bucket.Avatar
}

func (a *File) Upload(ctx context.Context, item schema.File) (*schema.IDResult, error) {
	info, err := a.AvatarModel.Upload(ctx, item.Name, item.Reader, item.Size, item.Type)
	return schema.NewIDResult(info.Key), err
}
