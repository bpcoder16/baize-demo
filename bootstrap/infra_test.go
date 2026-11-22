package bootstrap_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bpcoder16/baize-demo/bootstrap"
	"github.com/bpcoder16/baize-demo/internal/infra/aliyunes"
	"github.com/bpcoder16/baize-demo/internal/infra/aliyunoss"
	"github.com/bpcoder16/baize-demo/internal/infra/aliyunsts"
	"github.com/bpcoder16/baize-demo/internal/infra/lock"
	"github.com/bpcoder16/baize-demo/internal/infra/redis"
	"github.com/bpcoder16/baize-demo/internal/infra/sqlite"

	"github.com/alicebob/miniredis/v2"
	baize "github.com/bpcoder16/baize/bootstrap"
	baizeconstant "github.com/bpcoder16/baize/constant"
	"github.com/bpcoder16/pixiu/infra/env"
	"github.com/bpcoder16/pixiu/lifecycle"
	goredis "github.com/redis/go-redis/v9"
)

func TestInfraAccessors(t *testing.T) {
	// 命名注册表在关闭后不能重新打开，每次验证使用独立进程。
	if os.Getenv("BAIZE_DEMO_INFRA_TEST") == "" {
		cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestInfraAccessors$")
		cmd.Env = append(os.Environ(), "BAIZE_DEMO_INFRA_TEST=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v\n%s", err, out)
		}
		return
	}

	server := miniredis.RunT(t)
	dir := t.TempDir()
	config := fmt.Sprintf("host: %s\nport: %s\nmaxRetries: -1\n", server.Host(), server.Port())
	if err := os.WriteFile(filepath.Join(dir, "redis."+redis.DefaultResourceName+".yaml"), []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	plan := bootstrap.NewPlan()
	var resources lifecycle.Stack
	defer resources.Close()
	for _, resource := range plan.Resources {
		configDir := filepath.Join("..", "conf.example")
		switch resource.Name {
		case baizeconstant.ResourceRedis:
			configDir = dir
		case baizeconstant.ResourceSQLite, baizeconstant.ResourceAliyunOSS, baizeconstant.ResourceAliyunSTS, baizeconstant.ResourceAliyunES:
		default:
			continue
		}
		if err := resource.Init(configDir, &resources); err != nil {
			t.Fatal(err)
		}
	}
	ctx := t.Context()
	// 此用例只打开 SQLite、Redis 和云客户端，两个服务器数据库的迁移均显式关闭。
	if err := os.WriteFile(filepath.Join(dir, "mysql-migrate.default.yaml"), []byte("enabled: false"), 0600); err != nil {
		t.Fatal(err)
	}
	foundOutbox := false
	for _, step := range plan.Prepare {
		configDir := filepath.Join("..", "conf.example")
		if step.Name == "mysql-default-migrate" || step.Name == "jetstream-outbox" {
			configDir = dir
		}
		err := step.Init(ctx, configDir, &resources)
		if step.Name == "jetstream-outbox" {
			foundOutbox = true
			if err == nil || !strings.Contains(err.Error(), "migration enabled") {
				t.Fatalf("Outbox must reject disabled MySQL migration: %v", err)
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	if !foundOutbox {
		t.Fatal("Demo 缺少 JetStream Outbox Prepare 步骤")
	}

	db := sqlite.DefaultDB(ctx)
	var one int
	if err := db.Raw("SELECT 1").Scan(&one).Error; err != nil || one != 1 || db.Statement.Context != ctx {
		t.Fatalf("SQLite 会话查询结果 %d，错误 %v", one, err)
	}
	if aliyunoss.DefaultClient() == nil || aliyunsts.DefaultClient() == nil || aliyunes.DefaultClient() == nil {
		t.Fatal("未取得已装配的阿里云客户端")
	}
	client := redis.DefaultClient()
	if redis.DefaultInstance().Client() != client {
		t.Fatal("Redis 入口应共用同一连接池")
	}
	const key = "demo:test:job"
	err := lock.DefaultExecutor().Do(ctx, key, func(lockCtx context.Context) error {
		count, err := client.Exists(lockCtx, key).Result()
		if err != nil {
			return err
		}
		if count != 1 {
			return errors.New("在装配的 Redis 中未找到锁")
		}
		_, err = lock.DefaultExecutor().TryDo(lockCtx, key, func(context.Context) error {
			return errors.New("竞争者进入了临界区")
		})
		return err
	})
	if err != nil || server.Exists(key) {
		t.Fatalf("执行或释放锁失败，错误 %v", err)
	}
	t.Run("Cron模板绑定到宿主任务和共享Redis", func(t *testing.T) {
		if err := env.Init(env.Config{AppName: "demo", RunMode: "test", TimeLocation: "UTC", ConfigDirPath: dir, LocalIP: "127.0.0.1"}); err != nil {
			t.Fatal(err)
		}
		i := slices.IndexFunc(plan.Services, func(s baize.Service) bool { return s.Name == baizeconstant.ServiceCron })
		if i < 0 {
			t.Fatal("Demo 启动计划缺少 Cron 服务")
		}
		service, err := plan.Services[i].Build(t.Context(), filepath.Join("..", "conf.example"))
		if err != nil {
			t.Fatal(err)
		}
		if err := service.Stop(); err != nil {
			t.Fatal(err)
		}
		if err := service.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
	})
	if err := resources.Close(); err != nil {
		t.Fatal(err)
	}
	if err := client.Ping(ctx).Err(); !errors.Is(err, goredis.ErrClosed) {
		t.Fatalf("资源栈关闭后 Redis 返回 %v", err)
	}
}
