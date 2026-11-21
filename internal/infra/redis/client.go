package redis

import (
	"github.com/bpcoder16/pixiu/infra/redisx"
	goredis "github.com/redis/go-redis/v9"
)

// DefaultResourceName 是 Plan.Resources 注册的 Redis 实例名。
const DefaultResourceName = "default"

// DefaultInstance 返回装配好的 Pixiu Redis 实例。
func DefaultInstance() *redisx.Client {
	return redisx.Named(DefaultResourceName)
}

// DefaultClient 返回该实例共享的 Redis 客户端，执行命令时传入当前 ctx。
func DefaultClient() *goredis.Client {
	return DefaultInstance().Client()
}
