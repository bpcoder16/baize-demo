package aliyunoss

import "github.com/bpcoder16/pixiu/infra/aliyun/ossx"

// DefaultResourceName 是 Plan.Resources 注册的阿里云 OSS 实例名。
const DefaultResourceName = "default"

// DefaultClient 返回绑定示例 Bucket 的 OSS 客户端。
func DefaultClient() *ossx.Client {
	return ossx.Named(DefaultResourceName)
}
