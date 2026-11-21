// Package router 登记 Demo 的 Cron 任务。
package router

import (
	"github.com/bpcoder16/baize-demo/internal/cron/handler"
	"github.com/bpcoder16/pixiu/infra/cronx"
)

// Registry 按函数名登记任务，cron.yaml 的 jobs 按相同名称启用。
var Registry = cronx.MustNewHandlerRegistry(
	handler.Heartbeat,
)
