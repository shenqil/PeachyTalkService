package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// Config redis配置参数
type Config struct {
	Addr      string // 地址(IP:Port)
	DB        int    // 数据库
	Password  string // 密码
	KeyPrefix string // 储存key的前缀
}

// Store redis储存
type Store struct {
	cli    *redis.Client
	prefix string
}

// NewStore 创建基于redis储存的实例
func NewStore(cfg *Config) *Store {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return &Store{
		cli:    cli,
		prefix: cfg.KeyPrefix,
	}
}

// NewStoreWithClient 使用redis客户端创建储存实例
func NewStoreWithClient(cli *redis.Client, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

//// NewStoreWithClusterClient 使用redis集群客户端创建储存实例
//func NewStoreWithClusterClient(cli *redis.ClusterClient, keyPrefix string) *Store {
//	return &Store{
//		cli:    cli,
//		prefix: keyPrefix,
//	}
//}

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.prefix, key)
}

// Set ...
func (s *Store) Set(ctx context.Context, tokenString string, expiration time.Duration) error {
	cmd := s.cli.Set(ctx, s.wrapperKey(tokenString), "1", expiration)
	return cmd.Err()
}

// Delete ...
func (s *Store) Delete(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Del(ctx, s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() > 0, nil
}

// Check ...
func (s *Store) Check(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Exists(ctx, s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() > 0, nil
}

// Close ...
func (s *Store) Close() error {
	return s.cli.Close()
}
