package nats

import "github.com/bpcoder16/pixiu/infra/natsx"

// DefaultResourceName 是 Plan.Resources 注册的 NATS 实例名。
const DefaultResourceName = "default"

// DefaultClient 返回已装配的 NATS 客户端。
func DefaultClient() *natsx.Client {
	return natsx.Named(DefaultResourceName)
}
