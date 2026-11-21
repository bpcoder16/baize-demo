package clickhouse

import (
	"context"

	"github.com/bpcoder16/pixiu/infra/clickhousex"
	"gorm.io/gorm"
)

// DefaultResourceName 是 Plan.Resources 注册的 ClickHouse 实例名。
const DefaultResourceName = "default"

// DefaultMasterDB 返回绑定当前 ctx 的 ClickHouse 主库会话。
func DefaultMasterDB(ctx context.Context) *gorm.DB {
	return clickhousex.Named(DefaultResourceName).MasterDB(ctx)
}

// DefaultSlaveDB 返回绑定当前 ctx 的从库会话。未配置从库时使用主库。
func DefaultSlaveDB(ctx context.Context) *gorm.DB {
	return clickhousex.Named(DefaultResourceName).SlaveDB(ctx)
}
