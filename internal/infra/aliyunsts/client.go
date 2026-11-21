package aliyunsts

import "github.com/bpcoder16/pixiu/infra/aliyun/stsx"

// DefaultResourceName 是 Plan.Resources 注册的阿里云 STS 实例名。
const DefaultResourceName = "default"

// DefaultClient 返回用于申请 OSS 临时凭据的 STS 客户端。
func DefaultClient() *stsx.Client {
	return stsx.Named(DefaultResourceName)
}
