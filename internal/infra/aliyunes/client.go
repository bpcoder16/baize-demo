package aliyunes

import "github.com/bpcoder16/pixiu/infra/aliyun/esx"

// DefaultResourceName 是 Plan.Resources 注册的阿里云 ES 管理客户端实例名。
const DefaultResourceName = "default"

// DefaultClient 返回用于查询实例和管理词典的阿里云 ES 管理客户端。
func DefaultClient() *esx.Client {
	return esx.Named(DefaultResourceName)
}
