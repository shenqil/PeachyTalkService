package bucket

import (
	"context"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"io"
)

var UserFileSet = wire.NewSet(wire.Struct(new(UserFile), "*"))

var UserFileBucketName = "files"
var UserFileBucketLocation = "us-east-1"

type UserFile struct {
	minioClient *minio.Client
}

func (a *UserFile) Upload(ctx context.Context, fileName string, reader io.Reader, size int64, contentType string) (info minio.UploadInfo, err error) {
	return a.minioClient.PutObject(ctx, AvatarBucketName, fileName, reader, size, minio.PutObjectOptions{ContentType: contentType})
}
