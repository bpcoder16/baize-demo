package lock

import "github.com/bpcoder16/baize/service/lockx"

// DefaultPrepareExecutorName 是 Plan.Prepare 注册的 Redis 锁执行器名。
const DefaultPrepareExecutorName = "default"

// DefaultExecutor 返回 Prepare 阶段注册的 Redis 锁执行器。
func DefaultExecutor() *lockx.Executor {
	return lockx.Named(DefaultPrepareExecutorName)
}
