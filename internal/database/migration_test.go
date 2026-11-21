package database_test

import (
	"io/fs"
	"slices"
	"sync"
	"testing"

	"github.com/bpcoder16/baize-demo/internal/database"
	"github.com/bpcoder16/pixiu/infra/schemamigrate"
	"github.com/bpcoder16/pixiu/infra/sqlitex"
	"gorm.io/gorm/schema"
)

func TestMigrationSeedsDemoUser(t *testing.T) {
	columns := []string{
		"id", "user_number", "nickname", "phone", "avatar", "bio",
		"balance", "preferences", "remark",
		"is_disabled", "last_login_at", "created_at", "updated_at", "deleted_at",
	}
	mysqlSeeds := database.MySQLSeeds()
	postgresqlSeeds := database.PostgreSQLSeeds()
	tests := []struct {
		name      string
		tables    []schemamigrate.Table
		files     fs.FS
		seeds     []schemamigrate.Seed
		charset   string
		collation string
	}{
		{
			name:      "mysql",
			tables:    database.MySQLTables(),
			files:     mysqlSeeds.Files,
			seeds:     mysqlSeeds.Seeds,
			charset:   "utf8mb4",
			collation: "utf8mb4_unicode_ci",
		},
		{
			name:   "postgresql",
			tables: database.PostgreSQLTables(),
			files:  postgresqlSeeds.Files,
			seeds:  postgresqlSeeds.Seeds,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.tables) != 1 {
				t.Fatalf("迁移表数量为 %d，期望 1", len(tt.tables))
			}
			table := tt.tables[0]
			if table.Charset != tt.charset || table.Collation != tt.collation {
				t.Fatalf("字符集和排序规则为 %q / %q，期望 %q / %q",
					table.Charset, table.Collation, tt.charset, tt.collation)
			}
			model, err := schema.Parse(table.Model, &sync.Map{}, schema.NamingStrategy{})
			if err != nil {
				t.Fatal(err)
			}
			if model.Table != "demo_users" || len(model.PrimaryFields) != 1 ||
				model.PrimaryFields[0].DBName != "id" || !model.PrimaryFields[0].AutoIncrement {
				t.Fatal("用户示例应使用 demo_users 的自增 id 主键")
			}
			if !slices.Equal(model.DBNames, columns) {
				t.Fatalf("用户示例字段为 %v，期望 %v", model.DBNames, columns)
			}
			checkDemoUserSeeds(t, tt.files, tt.seeds)
		})
	}
}

func checkDemoUserSeeds(t *testing.T, files fs.FS, seeds []schemamigrate.Seed) {
	t.Helper()
	client, err := sqlitex.New(sqlitex.Config{Name: "test", DSN: ":memory:"})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := t.Context()
	// SQLite 只建种子用到的列，各方言的建表声明由真实数据库验证。
	db := client.DB(ctx)
	if err := db.Exec(`CREATE TABLE demo_users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_number TEXT NOT NULL UNIQUE,
		nickname TEXT NOT NULL
	)`).Error; err != nil {
		t.Fatal(err)
	}
	apply := func() {
		t.Helper()
		if err := schemamigrate.ApplySeeds(ctx, client.DB, files, seeds); err != nil {
			t.Fatal(err)
		}
	}
	readUser := func(number string) (int, string) {
		t.Helper()
		var id int
		var nickname string
		row := db.Raw("SELECT id, nickname FROM demo_users WHERE user_number = ?", number).Row()
		if err := row.Scan(&id, &nickname); err != nil {
			t.Fatal(err)
		}
		return id, nickname
	}

	apply()
	id, nickname := readUser("demo-user-001")
	if id <= 0 || nickname != "Demo User" {
		t.Fatalf("示例用户 ID 为 %d，昵称为 %q", id, nickname)
	}
	if err := db.Exec("UPDATE demo_users SET nickname = ? WHERE id = ?", "已有昵称", id).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("INSERT INTO demo_users (user_number, nickname) VALUES (?, ?)", "other", "Other").Error; err != nil {
		t.Fatal(err)
	}
	apply()
	if gotID, gotNickname := readUser("demo-user-001"); gotID != id || gotNickname != "已有昵称" {
		t.Fatalf("种子重跑后示例用户 ID 为 %d，昵称为 %q", gotID, gotNickname)
	}
	if _, nickname := readUser("other"); nickname != "Other" {
		t.Fatalf("种子重跑后其他用户昵称为 %q", nickname)
	}
}
