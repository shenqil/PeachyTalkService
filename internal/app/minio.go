package app

import (
	"PeachyTalkService/internal/app/config"
	"PeachyTalkService/internal/app/model/miniox"
	"github.com/minio/minio-go/v7"
)

func InitMinio() (*minio.Client, func(), error) {
	cfg := config.C.MINIO
	client, cleanFuc, err := NewMinioClient()
	if err != nil {
		return nil, cleanFuc, err
	}

	if cfg.AutoMakeBucket {
		err = miniox.AutoMakeBucket(client)
		if err != nil {
			return nil, cleanFuc, err
		}
	}

	return client, cleanFuc, nil
}

func NewMinioClient() (*minio.Client, func(), error) {
	cfg := config.C

	return miniox.NewMinioClient(&miniox.Config{
		Debug:           cfg.MINIO.Debug,
		Endpoint:        cfg.MINIO.Endpoint,
		AccessKeyID:     cfg.MINIO.AccessKeyID,
		SecretAccessKey: cfg.MINIO.SecretAccessKey,
		UseSSL:          cfg.MINIO.UseSSL,
	})
}
