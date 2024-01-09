package bucket

import (
	"context"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"io"
)

var AvatarSet = wire.NewSet(wire.Struct(new(Avatar), "*"))

var AvatarBucketName = "avatar"
var AvatarBucketLocation = "us-east-1"

type Avatar struct {
	MinioClient *minio.Client
}

func (a *Avatar) Upload(ctx context.Context, fileName string, reader io.Reader, size int64, contentType string) (info minio.UploadInfo, err error) {
	return a.MinioClient.PutObject(ctx, AvatarBucketName, fileName, reader, size, minio.PutObjectOptions{ContentType: contentType})
}
