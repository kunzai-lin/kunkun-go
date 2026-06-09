package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RDB 全局 Redis 客户端，在 InitRedis 成功后可用。
var RDB *redis.Client

// InitRedis 根据 config.yaml 中的 redis 段初始化连接。
func InitRedis() {
	addr := viper.GetString("redis.addr")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("redis ping failed: %w", err))
	}
}

// CloseRedis 在进程退出时关闭 Redis（可选，与 DB 对称）。
func CloseRedis() error {
	if RDB == nil {
		return nil
	}
	return RDB.Close()
}
