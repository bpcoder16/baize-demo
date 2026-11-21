package bootstrap_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bpcoder16/baize-demo/bootstrap"
	baize "github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/pixiu/infra/natsx"
	"github.com/nats-io/nats-server/v2/server"
)

// 分进程验证默认 Core 配置和取消注释后的 JetStream 示例，各自使用独立的命名客户端。
func TestDemoNATSConsumer(t *testing.T) {
	mode := os.Getenv("BAIZE_DEMO_NATS_TEST")
	if mode == "" {
		for _, mode := range []string{"core", "jetstream"} {
			t.Run(mode, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
				defer cancel()
				cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDemoNATSConsumer$")
				cmd.Env = append(os.Environ(), "BAIZE_DEMO_NATS_TEST="+mode)
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("%v\n%s", err, out)
				}
			})
		}
		return
	}
	services := bootstrap.NewPlan().Services
	idx := slices.IndexFunc(services, func(s baize.Service) bool { return s.Name == "nats-consumer:default" })
	if idx < 0 {
		t.Fatal("Demo 缺少 default 实例的 NATS 消费服务")
	}
	srv, err := server.NewServer(&server.Options{
		Host:      "127.0.0.1",
		Port:      -1,
		NoLog:     true,
		NoSigs:    true,
		JetStream: mode == "jetstream",
		StoreDir:  t.TempDir(),
	})
	if err != nil {
		t.Fatal(err)
	}
	srv.Start()
	defer func() {
		srv.Shutdown()
		srv.WaitForShutdown()
	}()
	if !srv.ReadyForConnections(5 * time.Second) {
		t.Fatal("临时 Core NATS 未就绪")
	}
	client, err := natsx.NewNamed(natsx.Config{Name: "default", URLs: []string{srv.ClientURL()}})
	if err != nil {
		t.Fatal(err)
	}
	defer natsx.CloseAll()
	dir := filepath.Join("..", "conf.example")
	if mode == "jetstream" {
		const filename = "nats-register.default.yaml"
		raw, readErr := os.ReadFile(filepath.Join(dir, filename))
		if readErr != nil {
			t.Fatal(readErr)
		}
		before, after, ok := strings.Cut(string(raw), "\n# streams:\n")
		if !ok {
			t.Fatal("模板缺少注释的 streams 示例")
		}
		block, rest, _ := strings.Cut(after, "\n\n")
		lines := strings.Split(block, "\n")
		for i := range lines {
			lines[i] = strings.TrimPrefix(lines[i], "# ")
		}
		enabled := before + "\nstreams:\n" + strings.Join(lines, "\n") + "\n\n" + rest
		dir = t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, filename), []byte(enabled), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	service, err := services[idx].Build(t.Context(), dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- service.Run() }()
	defer func() {
		service.Cancel()
		if err := service.Wait(ctx); err != nil {
			t.Error(err)
		}
	}()
	for !srv.GlobalAccount().SubscriptionInterest("demo.ping") {
		select {
		case err := <-done:
			t.Fatalf("服务提前退出: %v", err)
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(10 * time.Millisecond):
		}
	}
	reply, err := client.Request(ctx, "demo.ping", []byte("hello"))
	if err != nil || string(reply.Data) != "pong" {
		t.Fatalf("demo.ping 得到 %v，错误 %v", reply, err)
	}
	if mode == "jetstream" {
		if _, err := client.PublishJetStream(ctx, "demo.events.created", []byte(`{"id":"demo-1"}`), "demo-1"); err != nil {
			t.Fatal(err)
		}
		consumer, err := client.JetStream().Consumer(ctx, "DEMO_EVENTS", "demo-event-worker")
		if err != nil {
			t.Fatal(err)
		}
		for {
			info, infoErr := consumer.Info(ctx)
			if infoErr != nil {
				t.Fatal(infoErr)
			}
			if info.AckFloor.Stream == 1 {
				break
			}
			select {
			case <-ctx.Done():
				t.Fatal("示例 Handler 未确认消息")
			case <-time.After(10 * time.Millisecond):
			}
		}
	}
}
