package miniox

import (
	"PeachyTalkService/internal/app/model/miniox/bucket"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config 配置参数
type Config struct {
	Debug           bool
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

// NewMinioClient 创建一个minioClient实例
func NewMinioClient(c *Config) (*minio.Client, func(), error) {
	var minioClient *minio.Client
	var err error

	minioClient, err = minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure: c.UseSSL,
	})
	if err != nil {
		return nil, nil, err
	}

	cleanFunc := func() {}

	return minioClient, cleanFunc, nil
}

func makeBucket(ctx context.Context, minioClient *minio.Client, bucketName string, location string) error {
	// Make a new bucket called testbucket.
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists != nil && !exists {
			return err
		}
	}

	return nil
}

// AutoMakeBucket 自动创建存储桶
func AutoMakeBucket(minioClient *minio.Client) error {
	ctx := context.Background()
	var err error
	err = makeBucket(ctx, minioClient, bucket.AvatarBucketName, bucket.AvatarBucketLocation)
	if err != nil {
		return err
	}

	return nil
}
