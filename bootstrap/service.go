package bootstrap

import (
	"github.com/bpcoder16/baize-demo/internal/http/router"
	_ "github.com/bpcoder16/baize-demo/openapi" // 向 Swag 注册宿主的接口文档。

	"github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/httpboot"
	"github.com/bpcoder16/pixiu/infra/ginx"
)

func services() []bootstrap.Service {
	return []bootstrap.Service{
		httpboot.Gin(
			router.Register,
			ginx.OpenAPIModule{Name: "app", Title: "App"},
			ginx.OpenAPIModule{Name: "admin", Title: "Admin"},
		),
	}
}
