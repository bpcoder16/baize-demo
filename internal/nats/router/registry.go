// Package router 登记 Demo 的 NATS 消费函数。
package router

import (
	"github.com/bpcoder16/baize-demo/internal/nats/handler"
	"github.com/bpcoder16/pixiu/infra/natsx"
)

// JetStreamRegistry 按函数名绑定 streams[].consumers 中的 handler。
var JetStreamRegistry = natsx.MustNewJetStreamHandlerRegistry(
	handler.HandleDemoEvent,
)

// CoreRegistry 按函数名绑定 coreConsumers 中的 handler。
var CoreRegistry = natsx.MustNewCoreHandlerRegistry(
	handler.HandlePing,
)
