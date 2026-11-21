package database

import (
	"github.com/bpcoder16/baize-demo/internal/model/entity/defaultmysql"
	"github.com/bpcoder16/pixiu/infra/schemamigrate"
)

// MySQLTables 返回 MySQL 默认实例的迁移表清单。
func MySQLTables() []schemamigrate.Table {
	return []schemamigrate.Table{
		{
			Model:     &defaultmysql.DemoUser{},
			Comment:   "示例用户表",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
		},
	}
}
