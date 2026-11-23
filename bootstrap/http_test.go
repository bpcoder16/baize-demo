package bootstrap_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/bpcoder16/baize-demo/bootstrap"
	"github.com/bpcoder16/baize-demo/internal/http/router"
	baize "github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/httpboot"
	"github.com/bpcoder16/pixiu/infra/env"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

func TestDemoHTTPModules(t *testing.T) {
	// Baize 构建 HTTP 服务时读取进程级环境，在独立进程中加载样例配置。
	if os.Getenv("BAIZE_DEMO_TEST_HTTP") == "" {
		cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestDemoHTTPModules$")
		cmd.Env = append(os.Environ(), "BAIZE_DEMO_TEST_HTTP=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v\n%s", err, output)
		}
		return
	}
	plan := bootstrap.NewPlan()
	if !slices.ContainsFunc(plan.Services, func(s baize.Service) bool { return s.Name == "http" }) {
		t.Fatal("Demo 启动计划缺少 HTTP 服务")
	}
	dir, err := filepath.Abs("../conf.example")
	if err != nil {
		t.Fatal(err)
	}
	if err := env.Init(env.Config{
		AppName:       "baize-demo",
		RunMode:       env.RunModeRelease,
		TimeLocation:  "UTC",
		ConfigDirPath: dir,
		LocalIP:       "127.0.0.1",
	}); err != nil {
		t.Fatal(err)
	}
	var r *gin.Engine
	service := httpboot.Gin(func(engine *gin.Engine) {
		router.Register(engine)
		r = engine
	})
	if _, err := service.Build(t.Context(), dir); err != nil {
		t.Fatal(err)
	}
	const origin = "http://localhost:5173"
	preflight := httptest.NewRequest(http.MethodOptions, "/app/ping", nil)
	preflight.Header.Set("Origin", origin)
	preflight.Header.Set("Access-Control-Request-Method", http.MethodGet)
	preflight.Header.Set("Access-Control-Request-Headers", "Authorization")
	response := httptest.NewRecorder()
	r.ServeHTTP(response, preflight)
	if response.Code != http.StatusNoContent || response.Header().Get("Access-Control-Allow-Origin") != origin {
		t.Fatalf("App 预检得到 %d %v", response.Code, response.Header())
	}
	if _, err := swag.ReadDoc("web"); err == nil {
		t.Error("Demo 仍注册了 Web 文档")
	}
	web := httptest.NewRecorder()
	r.ServeHTTP(web, httptest.NewRequest(http.MethodGet, "/web/ping", nil))
	if web.Code != http.StatusNotFound {
		t.Errorf("GET /web/ping = %d，预期 404", web.Code)
	}
	modules := []string{"app", "admin"}
	for _, name := range modules {
		path := "/" + name + "/ping"
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Origin", origin)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK || w.Body.String() != "pong" {
			t.Fatalf("GET %s = %d %q", path, w.Code, w.Body.String())
		}
		wantOrigin := ""
		if name == "app" {
			wantOrigin = origin
		}
		if got := w.Header().Get("Access-Control-Allow-Origin"); got != wantOrigin {
			t.Errorf("GET %s 跨域来源为 %q，期望 %q", path, got, wantOrigin)
		}
		raw, err := swag.ReadDoc(name)
		if err != nil {
			t.Fatal(err)
		}
		var doc struct {
			Paths map[string]map[string]json.RawMessage
		}
		if err := json.Unmarshal([]byte(raw), &doc); err != nil {
			t.Fatal(err)
		}
		if doc.Paths[path]["get"] == nil {
			t.Errorf("%s 文档缺少 GET %s", name, path)
		}
		for _, other := range modules {
			if other != name && doc.Paths["/"+other+"/ping"] != nil {
				t.Errorf("%s 文档中出现了 %s 的示例接口", name, other)
			}
		}
	}
}
