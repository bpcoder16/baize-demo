package bootstrap

import (
	"github.com/bpcoder16/baize-demo/internal/infra/aliyunes"
	"github.com/bpcoder16/baize-demo/internal/infra/aliyunoss"
	"github.com/bpcoder16/baize-demo/internal/infra/aliyunsts"
	"github.com/bpcoder16/baize-demo/internal/infra/clickhouse"
	"github.com/bpcoder16/baize-demo/internal/infra/elasticsearch"
	"github.com/bpcoder16/baize-demo/internal/infra/mysql"
	"github.com/bpcoder16/baize-demo/internal/infra/nats"
	"github.com/bpcoder16/baize-demo/internal/infra/postgresql"
	"github.com/bpcoder16/baize-demo/internal/infra/redis"
	"github.com/bpcoder16/baize-demo/internal/infra/sqlite"

	baize "github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/aliyunesboot"
	"github.com/bpcoder16/baize/bootstrap/aliyunossboot"
	"github.com/bpcoder16/baize/bootstrap/aliyunstsboot"
	"github.com/bpcoder16/baize/bootstrap/clickhouseboot"
	esv8 "github.com/bpcoder16/baize/bootstrap/elasticsearchboot/v8"
	"github.com/bpcoder16/baize/bootstrap/mysqlboot"
	"github.com/bpcoder16/baize/bootstrap/natsboot"
	"github.com/bpcoder16/baize/bootstrap/postgresqlboot"
	"github.com/bpcoder16/baize/bootstrap/redisboot"
	"github.com/bpcoder16/baize/bootstrap/sqliteboot"
)

func resources() []baize.Resource {
	return []baize.Resource{
		mysqlboot.Instances(mysql.DefaultResourceName),
		postgresqlboot.Instances(postgresql.DefaultResourceName),
		clickhouseboot.Instances(clickhouse.DefaultResourceName),
		sqliteboot.Instances(sqlite.DefaultResourceName),
		redisboot.Instances(redis.DefaultResourceName),
		natsboot.Instances(nats.DefaultResourceName),
		esv8.Instances(elasticsearch.DefaultResourceName),
		aliyunossboot.Instances(aliyunoss.DefaultResourceName),
		aliyunstsboot.Instances(aliyunsts.DefaultResourceName),
		aliyunesboot.Instances(aliyunes.DefaultResourceName),
	}
}
