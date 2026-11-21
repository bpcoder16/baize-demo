package bootstrap_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/bpcoder16/baize-demo/bootstrap"
	"github.com/bpcoder16/baize-demo/internal/http/router"
	baize "github.com/bpcoder16/baize/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

func TestDemoHTTPModules(t *testing.T) {
	plan := bootstrap.NewPlan()
	if !slices.ContainsFunc(plan.Services, func(s baize.Service) bool { return s.Name == "http" }) {
		t.Fatal("Demo 启动计划缺少 HTTP 服务")
	}
	r := gin.New()
	router.Register(r)
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
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != http.StatusOK || w.Body.String() != "pong" {
			t.Fatalf("GET %s = %d %q", path, w.Code, w.Body.String())
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
