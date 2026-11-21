package elasticsearch

import "github.com/bpcoder16/pixiu/infra/elasticsearchx"

// DefaultResourceName 是 Plan.Resources 注册的 Elasticsearch 实例名。
const DefaultResourceName = "default"

// DefaultClient 返回用于索引和文档操作的 Elasticsearch 客户端。
func DefaultClient() *elasticsearchx.Client {
	return elasticsearchx.Named(DefaultResourceName)
}
