package bootstrap

import (
	cronrouter "github.com/bpcoder16/baize-demo/internal/cron/router"
	httprouter "github.com/bpcoder16/baize-demo/internal/http/router"
	"github.com/bpcoder16/baize-demo/internal/infra/nats"
	"github.com/bpcoder16/baize-demo/internal/infra/realtime"
	natsrouter "github.com/bpcoder16/baize-demo/internal/nats/router"
	realtimerouter "github.com/bpcoder16/baize-demo/internal/realtime/router"
	_ "github.com/bpcoder16/baize-demo/openapi" // 向 Swag 注册宿主的接口文档。

	"github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/bootstrap/cronboot"
	"github.com/bpcoder16/baize/bootstrap/httpboot"
	"github.com/bpcoder16/baize/bootstrap/natsboot"
	"github.com/bpcoder16/baize/bootstrap/realtimeboot"
	"github.com/bpcoder16/pixiu/infra/ginx"
)

func services() []bootstrap.Service {
	return []bootstrap.Service{
		natsboot.ConsumerService(
			nats.DefaultResourceName,
			natsrouter.JetStreamRegistry,
			natsrouter.CoreRegistry,
		),
		cronboot.Service(cronrouter.Registry),
		httpboot.Gin(
			httprouter.Register,
			ginx.OpenAPIModule{Name: "app", Title: "App"},
			ginx.OpenAPIModule{Name: "admin", Title: "Admin"},
		),
		realtimeboot.Service(realtime.AppName, realtimerouter.NewAppServer),
	}
}
