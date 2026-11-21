package database

import (
	"embed"

	"github.com/bpcoder16/baize-demo/internal/model/entity/defaultpostgresql"
	"github.com/bpcoder16/baize/bootstrap/postgresqlboot"
	"github.com/bpcoder16/pixiu/infra/schemamigrate"
)

//go:embed seeds/users.sql
var postgresqlSeedFiles embed.FS

// PostgreSQLSeeds 返回 PostgreSQL 默认实例的种子选项。
func PostgreSQLSeeds() postgresqlboot.SeedOptions {
	return postgresqlboot.SeedOptions{
		Files: postgresqlSeedFiles,
		Seeds: []schemamigrate.Seed{
			{
				Model: &defaultpostgresql.DemoUser{},
				Where: map[string]any{"user_number": "demo-user-001"},
				File:  "seeds/users.sql",
			},
		},
	}
}
