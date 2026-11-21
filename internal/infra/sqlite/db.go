package sqlite

import (
	"context"

	"github.com/bpcoder16/pixiu/infra/sqlitex"
	"gorm.io/gorm"
)

// DefaultResourceName 是 Plan.Resources 注册的 SQLite 实例名。
const DefaultResourceName = "default"

// DefaultDB 返回绑定当前 ctx 的 SQLite 会话。
func DefaultDB(ctx context.Context) *gorm.DB {
	return sqlitex.Named(DefaultResourceName).DB(ctx)
}
