package bootstrap

import (
	"time"

	"github.com/bpcoder16/baize-demo/internal/infra/lock"
	"github.com/bpcoder16/baize-demo/internal/infra/redis"

	baize "github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/redisboot"
	"github.com/bpcoder16/pixiu/biz/lockx/redislock"
)

func prepareSteps() []baize.Step {
	return []baize.Step{
		redisboot.PrepareLock(lock.DefaultPrepareExecutorName, redis.DefaultResourceName, redislock.Config{
			TTL:           time.Minute,           // 锁的有效期，到期自动失效
			WaitTimeout:   5 * time.Second,       // 阻塞获取锁的最长等待时间
			RetryInterval: 50 * time.Millisecond, // 阻塞获取锁时，竞争失败后的重试间隔基准
		}),
	}
}
