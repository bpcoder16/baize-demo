package mysql

import (
	"context"

	"github.com/bpcoder16/pixiu/infra/mysqlx"
	"gorm.io/gorm"
)

// DefaultResourceName 是 Plan.Resources 注册的 MySQL 实例名。
const DefaultResourceName = "default"

// DefaultMasterDB 返回绑定当前 ctx 的 MySQL 主库会话。
func DefaultMasterDB(ctx context.Context) *gorm.DB {
	return mysqlx.Named(DefaultResourceName).MasterDB(ctx)
}

// DefaultSlaveDB 返回绑定当前 ctx 的从库会话。未配置从库时使用主库。
func DefaultSlaveDB(ctx context.Context) *gorm.DB {
	return mysqlx.Named(DefaultResourceName).SlaveDB(ctx)
}
