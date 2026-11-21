package database

import (
	"github.com/bpcoder16/baize-demo/internal/model/entity/defaultpostgresql"
	"github.com/bpcoder16/pixiu/infra/schemamigrate"
)

// PostgreSQLTables 返回 PostgreSQL 默认实例的迁移表清单。
func PostgreSQLTables() []schemamigrate.Table {
	return []schemamigrate.Table{
		{
			Model:   &defaultpostgresql.DemoUser{},
			Comment: "示例用户表",
		},
	}
}
