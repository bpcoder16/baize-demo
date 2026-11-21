// Package router 装配 Demo 的 Realtime 服务并登记上行消息路由。
package router

import (
	"context"
	"net/http"

	"github.com/bpcoder16/baize-demo/internal/realtime/handler"
	"github.com/bpcoder16/baize/service/realtime"
	"github.com/bpcoder16/pixiu/infra/websocketx"
)

// NewAppServer 为本机 ping 示例显式允许匿名连接。
func NewAppServer() (*realtime.Server, error) {
	return realtime.NewServer(realtime.Config{
		Transport: websocketx.Config{
			ReadLimit:    64 << 10, // 单条上行消息最多 64 KiB。
			MaxSendBytes: 64 << 10, // 单条下行消息最多 64 KiB。
		},
		MaxTopicsPerSession: 500, // 每个连接跨所有命名空间的订阅总数上限。
		Authenticate: func(context.Context, *http.Request) (realtime.ClientInfo, error) {
			return realtime.ClientInfo{}, nil
		},
		RegisterHandlers: Register,
	})
}

func Register(r *realtime.Router) error {
	return r.Handle("demo", "ping", handler.Ping)
}
