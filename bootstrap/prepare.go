package bootstrap

import (
	"time"

	"github.com/bpcoder16/baize-demo/internal/database"
	"github.com/bpcoder16/baize-demo/internal/infra/lock"
	"github.com/bpcoder16/baize-demo/internal/infra/mysql"
	"github.com/bpcoder16/baize-demo/internal/infra/nats"
	"github.com/bpcoder16/baize-demo/internal/infra/postgresql"
	"github.com/bpcoder16/baize-demo/internal/infra/redis"

	baize "github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/mysqlboot"
	"github.com/bpcoder16/baize/bootstrap/outboxboot"
	"github.com/bpcoder16/baize/bootstrap/postgresqlboot"
	"github.com/bpcoder16/baize/bootstrap/redisboot"
	"github.com/bpcoder16/baize/constant"
	"github.com/bpcoder16/pixiu/biz/lockx/redislock"
)

func prepareSteps() []baize.Step {
	return []baize.Step{
		mysqlboot.PrepareMigration(
			mysql.DefaultResourceName,
			database.MySQLTables(),
			database.MySQLSeeds(),
		),
		postgresqlboot.PrepareMigration(
			postgresql.DefaultResourceName,
			database.PostgreSQLTables(),
			database.PostgreSQLSeeds(),
		),
		redisboot.PrepareLock(lock.DefaultPrepareExecutorName, redis.DefaultResourceName, redislock.Config{
			TTL:           time.Minute,           // 锁的有效期，到期自动失效
			WaitTimeout:   5 * time.Second,       // 阻塞获取锁的最长等待时间
			RetryInterval: 50 * time.Millisecond, // 阻塞获取锁时，竞争失败后的重试间隔基准
		}),
		outboxboot.PrepareJetStream(outboxboot.JetStreamConfig{
			Backend:         constant.MySQL,
			DatabaseName:    mysql.DefaultResourceName,
			ResourceNames:   []string{nats.DefaultResourceName},
			CronTaskEvery:   10 * time.Second,
			CronTaskTimeout: time.Minute,
		}),
	}
}
