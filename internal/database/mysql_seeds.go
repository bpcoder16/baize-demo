// Package database 声明 demo 自己维护的表和种子。
package database

import (
	"embed"

	"github.com/bpcoder16/baize-demo/internal/model/entity/defaultmysql"
	"github.com/bpcoder16/baize/bootstrap/mysqlboot"
	"github.com/bpcoder16/pixiu/infra/schemamigrate"
)

//go:embed seeds/users.sql
var mysqlSeedFiles embed.FS

// MySQLSeeds 返回 MySQL 默认实例的种子选项。
func MySQLSeeds() mysqlboot.SeedOptions {
	return mysqlboot.SeedOptions{
		Files: mysqlSeedFiles,
		Seeds: []schemamigrate.Seed{
			{
				Model: &defaultmysql.DemoUser{},
				Where: map[string]any{"user_number": "demo-user-001"},
				File:  "seeds/users.sql",
			},
		},
	}
}
