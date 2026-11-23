# Baize Demo

运行 `go run .` 会初始化环境、时区和日志，打开 MySQL、PostgreSQL、ClickHouse、SQLite、Redis、NATS 和 Elasticsearch，再在本地构造阿里云 OSS、STS 和 ES 管理端客户端。准备步骤完成后，进程启动 Realtime WebSocket、Gin HTTP、Cron 和 `default` 实例的 NATS 消费服务。Demo 使用 Baize 的 Redis 锁和 JetStream Outbox，依赖版本见 [go.mod](go.mod)。运行需要 Go 1.27.3、可用的 C 编译器，以及配置中指定的数据库和网络服务。SQLite 默认使用进程内存。

各组件在 `internal/infra/<组件>/` 中维护自己的资源实例名。Redis 锁的装配引用 [lock.DefaultPrepareExecutorName](internal/infra/lock/executor.go) 和 [redis.DefaultResourceName](internal/infra/redis/client.go)。[bootstrap/resource.go](bootstrap/resource.go) 声明共享资源，[bootstrap/prepare.go](bootstrap/prepare.go) 声明准备步骤。[bootstrap/service.go](bootstrap/service.go) 声明 Realtime WebSocket、Gin HTTP、Cron 和 NATS 消费服务。[NewPlan](bootstrap/plan.go) 将三份清单组装为启动计划，由根入口 `main.go` 交给白泽的 `RunCLI` 执行。`Plan.Prepare` 中的 `redisboot.PrepareLock` 用名为 `default` 的 Redis 创建同名锁执行器。业务代码从 `github.com/bpcoder16/baize-demo/internal/infra/<组件>` 导入对应包，在资源或准备步骤就绪后获取共享句柄。客户端由白泽统一关闭。首次运行或升级时，从模板补齐本地配置。

```sh
mkdir -p conf
cp -n conf.example/*.yaml conf/
```

首次运行 `go run .` 或执行 `go test ./...` 前，先在仓库根目录生成接口文档。`openapi/app.go` 和 `openapi/admin.go` 声明各模块的文档信息，生成指令集中在 `openapi/placeholder.go`。Git 仅放行这三个源文件。

```sh
go generate ./openapi
```

应用配置在 `conf/app.yaml`。应用和实例配置中的超时建议写明单位，例如 `"30s"` 或 `"500ms"`。YAML 整数按纳秒解析，`resourceCloseTimeout: 30` 表示 30 纳秒。共享资源各有一份连接配置。已有部署升级时，将 `aliyunesx.search-control.yaml` 和 `aliyunstsx.oss.yaml` 分别改名为 `aliyunesx.default.yaml` 和 `aliyunstsx.default.yaml`，沿用原有配置内容。

- `conf/mysql.default.yaml` 填写数据库地址、库名和凭据。业务通过 `mysql.DefaultMasterDB(ctx)` 和 `mysql.DefaultSlaveDB(ctx)` 获取会话。
- `conf/postgresql.default.yaml` 填写主库、可选从库及凭据。业务通过 `postgresql.DefaultMasterDB(ctx)` 和 `postgresql.DefaultSlaveDB(ctx)` 获取会话。模板使用 `sslMode: "verify-full"`，本地配置按实际服务填写 TLS 模式和证书。`sslRootCert`、`sslCert` 和 `sslKey` 的相对路径从进程工作目录解析。例如，从 `/srv/demo` 使用 `--conf /etc/demo/app.yaml` 启动时，`certs/ca.pem` 指向 `/srv/demo/certs/ca.pem`。证书随配置部署时可写绝对路径，如 `/etc/demo/certs/ca.pem`。
- `conf/clickhouse.default.yaml` 填写库名、凭据、协议和地址。Native 默认端口 9000，HTTP 默认端口 8123。业务通过 `clickhouse.DefaultMasterDB(ctx)` 和 `clickhouse.DefaultSlaveDB(ctx)` 获取会话。
- `conf/sqlite.default.yaml` 默认使用 `:memory:`，业务通过 `sqlite.DefaultDB(ctx)` 获取会话。持久化时将 `dsn` 改为文件路径或 `file:` URI，并提前创建父目录。构建 demo 时使用 `CGO_ENABLED=1`。
- `conf/redis.default.yaml` 填写主机、端口、认证信息和数据库编号。通过 `redis.DefaultClient()` 执行命令，需要 Pixiu 包装类型时使用 `redis.DefaultInstance()`。锁复用该实例的连接池，需要保留 `maxRetries: -1` 以关闭 Redis 命令重试。
- `conf/nats.default.yaml` 填写集群节点地址，启用认证时填写用户名和密码。默认启动会初始化 JetStream Outbox，因此 NATS 服务端须开启 JetStream。业务通过 `nats.DefaultClient()` 发布消息。
- `conf/elasticsearch.default.yaml` 填写节点地址和认证信息。启动计划使用 `elasticsearchboot/v8`，启动时通过 `GET /` 校验集群为 Elasticsearch 8。升级旧配置时删除 `version` 字段。使用自建 CA 时填写 `caCertFile`，相对路径从该文件所在目录解析。
- `conf/aliyunossx.default.yaml` 配置地域、Bucket 和访问密钥，业务通过 `aliyunoss.DefaultClient()` 访问该桶。
- `conf/aliyunstsx.default.yaml` 配置 STS 地域和访问密钥，业务通过 `aliyunsts.DefaultClient()` 申请临时凭据。
- `conf/aliyunesx.default.yaml` 填写阿里云 ES 实例的地域、`instanceID` 和管理 API 访问密钥。

这三个阿里云客户端都在首次业务调用时才访问云服务，模板中的示例值可用于启动 demo。实际调用时，在本地配置中填写云资源信息和获授权的凭据。收到 SIGINT 或 SIGTERM 后，Demo 先停止接收 WebSocket 连接和消息、HTTP 请求、NATS 消息和调度新任务，等待在途工作结束，再按声明的逆序回收客户端，最后关闭日志。所有资源共享 `conf/app.yaml` 的 `resourceCloseTimeout`（默认 30 秒），单个 NATS 客户端的 `closeTimeout`（默认 20 秒）计入这段时间。

### Realtime WebSocket

`realtime:app` 读取 `conf/realtime.app.yaml`，默认在本机 `ws://127.0.0.1:8889/ws/app` 接受匿名连接。`GET /ready/ws` 在同一端口返回 204 时可连接。路由在 `internal/realtime/router` 登记，消息由 `internal/realtime/handler` 处理。业务可通过 `internal/infra/realtime.App()` 取得命名实例。将 `"realtime:app"` 加入 `disabledServices` 可停用该实例。[服务配置](conf.example/realtime.app.yaml)

客户端选择子协议 `baize.realtime.v1.json` 后发送以下文本帧。服务端回复 `demo.pong` 和原请求 ID。选择 `baize.realtime.v1.protobuf` 时使用二进制 Envelope，ping 和 pong 的正文分别为 `google.protobuf.Empty` 和 `google.protobuf.StringValue`。

```json
{"kind":1,"requestId":"r1","scene":"demo","action":"ping","data":{}}
```

```json
{"kind":2,"requestId":"r1","type":"demo.pong","data":"pong"}
```

### Gin HTTP 与接口文档

首次运行或升级时，将 `conf.example/gin.yaml` 和 `conf.example/http.yaml` 复制到入口配置同目录。示例监听 `127.0.0.1:8888`，可在 `http.yaml` 的 `addr` 修改。启动 Demo 后，请求 `GET /app/ping` 或 `GET /admin/ping`，均返回 `pong`。

在 debug/test 模式下打开 [本地接口文档](http://127.0.0.1:8888/docs)，首页展示 App，可切换到 Admin。各模块的页面位于 `/docs/<模块名>`，Swagger 2.0 JSON 位于 `/openapi/<模块名>.json`。release 模式只注册业务路由。

`internal/http/router` 按 App、Admin 分文件挂载路由。接口实现和 Swag 注释放在 `internal/http/handler`。两端共用 Handler 时，在函数旁各写一个独立注释块，分别声明模块标记和完整的 `@Router` 路径。`@Tags` 用于文档内分组。修改注释后，在仓库根目录执行 `go generate ./openapi` 更新两份文档。[路由入口](internal/http/router/http.go) · [Ping 接口与文档注释](internal/http/handler/ping.go)

### Cron 定时任务

`internal/cron/router` 登记任务，`internal/cron/handler` 实现任务。启动时，`bootstrap/service.go` 将注册表交给 `cronboot.Service`。Demo 启动一分钟后首次执行 `Heartbeat`，随后每分钟一次。可在 `log/baize-demo/biz/cron.info.log` 中查看该任务的 `CronJobResult`。[任务注册表](internal/cron/router/registry.go) · [Heartbeat 实现](internal/cron/handler/heartbeat.go) · [Cron 模板](conf.example/cron.yaml)

`conf/cron.yaml` 中的任务名与函数名大小写一致。每项选用 `every`、`cron` 或 `randomEvery` 中的一种调度方式。示例借用 `default` Redis，通过 `baize-demo:cron:` 前缀协调同一任务的多个副本。单实例部署可省略 `lock`。修改配置后重启生效。

设置 `jobs: []` 会停用示例中的 `Heartbeat`，由 `bootstrap/prepare.go` 登记的 Outbox 补偿任务 `baize.outbox.retry.nats_jetstream` 仍会运行。当前 Demo 需要 Cron 来执行补偿，如果在 `app.yaml` 的 `disabledServices` 中加入 `"cron"`，Outbox 准备步骤会报错并中止启动。

### NATS 消费与消息应答

`internal/nats/router` 分别登记 `JetStreamRegistry` 和 `CoreRegistry`，`internal/nats/handler` 实现消费函数。`nats-consumer:default` 借用现有的 `default` 连接，读取同目录的 `nats-register.default.yaml`。配置中的 `handler` 与对应注册表的名称大小写一致。[注册表](internal/nats/router/registry.go) · [订阅配置](conf.example/nats-register.default.yaml)

启动时先绑定 Handler 并校验整份订阅配置，再注册远端资源。默认订阅配置只启用 Core 广播请求应答。需要多个实例分担请求时，取消 `queue` 的注释即可。每个实例默认建立一个同组订阅，增加并发时可将 `subscribers` 设为更大的正整数。启动 Demo 后，在连接同一 NATS 服务的 CLI 中执行下列命令，预期收到 `pong`。[Ping 实现](internal/nats/handler/ping.go)

```sh
nats request demo.ping hello
```

要启用 JetStream 订阅，取消 `conf/nats-register.default.yaml` 中 `streams` 块的首层注释并重启 Demo。运行下列命令后，`HandleDemoEvent` 会把事件 ID 写入本次消费日志，由 natsx 确认成功处理的消息。日志位于 `log/baize-demo/biz/nats.info.log`。命令使用模板默认的空 `namespace`，若设置了命名空间，CLI 主题需加上相应前缀。[事件处理函数](internal/nats/handler/event.go)

```sh
nats pub demo.events.created '{"id":"demo-1"}'
```

在 `app.yaml` 的 `disabledServices` 中加入 `"nats-consumer:default"` 可按实例停用消费，已声明的 NATS 连接仍可发布消息。

### JetStream Outbox

`bootstrap/prepare.go` 通过 `outboxboot.JetStreamConfig` 选择 MySQL 主库、发布实例、`CronTaskEvery` 和 `CronTaskTimeout`。示例运行 `baize.outbox.retry.nats_jetstream`，每 10 秒触发，单轮超时 1 分钟。启动需要启用 Cron，发布目标 subject 须有对应的 Stream。

运行下面的示例前，取消 `conf/nats-register.default.yaml` 中 `streams` 块的首层注释并启动 Demo，让 `DEMO_EVENTS` 接收 `demo.events.created`，由 `HandleDemoEvent` 处理。[订阅模板](conf.example/nats-register.default.yaml)

业务通过 `github.com/bpcoder16/baize/service/outbox` 在写业务数据的同一个默认 MySQL 主库事务中登记任务，并把入队错误返回给事务调用方。

```go
var id int64
err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&order).Error; err != nil {
        return err
    }
    var err error
    id, err = outbox.EnqueueJetStreamTx(tx, outbox.JetStreamMessage{
        Subject:   "demo.events.created",
        Data:      []byte(`{"id":"demo-1"}`),
        Header:    map[string][]string{},
        ExpiresAt: time.Now().Add(24 * time.Hour),
    })
    return err
})
if err != nil {
    return err
}
outbox.SendAsyncJetStream(id)
return nil
```

`ResourceName` 省略时使用 Prepare 中的第一个 NATS 实例，指定时须在名单内。登记时持久化实际实例和 `nats_jetstream` 类型。最外层事务提交后，也可调用 `outbox.SendJetStream(ctx, id)` 等待结果。补偿按该类型执行，同表同类型的副本使用一致的发布名单。

### 数据库迁移

`bootstrap/prepare.go` 通过 Baize 的 `mysqlboot.PrepareMigration` 和 `postgresqlboot.PrepareMigration` 为两个 `default` 实例声明迁移步骤。首次运行或升级时，将下面两份模板复制到 `app.yaml` 同目录。示例开启 MySQL 迁移以准备 Outbox 和业务表，PostgreSQL 迁移默认关闭。[MySQL 迁移配置](conf.example/mysql-migrate.default.yaml) · [PostgreSQL 迁移配置](conf.example/postgresql-migrate.default.yaml)

两库均创建或更新 `demo_users`。MySQL 使用 `utf8mb4` 字符集、`utf8mb4_unicode_ci` 排序规则和自增主键。PostgreSQL 使用 `bigint GENERATED BY DEFAULT AS IDENTITY` 主键、布尔禁用状态和 `timestamptz(3)` 时间列。[MySQL 模型](internal/model/entity/defaultmysql/demo.go) · [PostgreSQL 模型](internal/model/entity/defaultpostgresql/demo.go)

两库共用种子 SQL，按 `user_number = 'demo-user-001'` 校准。首次插入由数据库生成 ID，已有对应记录时保留现有值。旧示例表 `demo_statuses` 由部署方在目标数据库中手动清理。[用户种子](internal/database/seeds/users.sql) · [MySQL 种子选项](internal/database/mysql_seeds.go) · [PostgreSQL 种子选项](internal/database/postgresql_seeds.go)

每个实例的迁移配置用 `timeout` 设置迁表和种子共用的总时限，默认 `5m`。多副本部署时，由一个进程在 `mysql-migrate.default.yaml` 设置 `enabled: true`，迁移业务表、校准种子，并在数据库会话锁内准备 `baize_outbox`。完成后，普通副本设为 `enabled: false` 即可使用已有表。缺表时 Outbox 在 Prepare 阶段报错。

### Redis 锁

`bootstrap/prepare.go` 显式填写了 Baize 的默认锁策略。租期 `TTL` 为 1 分钟，阻塞获取锁最多等待 `WaitTimeout`（5 秒），竞争后的重试间隔以 `RetryInterval`（50 毫秒）为基准。业务代码通过 `internal/infra/lock` 的 `DefaultExecutor()` 获取该执行器。以下示例以任务 ID 作为锁键。

```go
func runJob(ctx context.Context, jobID string, work func(context.Context) error) error {
    return lock.DefaultExecutor().Do(ctx, "demo:job:"+jobID, work)
}
```

### Elasticsearch 同义词验证

词库文件部署到 ES 节点后，业务可通过 Demo 已装配的 `elasticsearch.DefaultClient()` 验证文件。选项类型来自 `github.com/bpcoder16/pixiu/infra/elasticsearchx`，路径相对于 ES 节点的配置目录。

```go
func validateSynonyms(ctx context.Context) error {
    result, err := elasticsearch.DefaultClient().ValidateSynonyms(ctx, elasticsearchx.SynonymValidationOptions{
        SynonymsPath: "analysis/synonyms.txt",
    })
    if err != nil && result.MayExist {
        return fmt.Errorf("请核查临时索引 %s: %w", result.Index, err)
    }
    return err
}
```

### 阿里云 ES 管理端

通过 `aliyunes.DefaultClient()` 查询云实例状态。索引和文档读写使用 `elasticsearch.DefaultClient()`。

```go
func instanceStatus(ctx context.Context) (string, error) {
    result, err := aliyunes.DefaultClient().DescribeInstance(ctx)
    if err != nil {
        return "", err
    }
    return result.Status, nil
}
```

```sh
go run .
cat log/baize-demo/biz/info.log
```

可用 `go run . --conf /etc/cards/staging.yaml` 指定应用配置（简写为 `-c`），各实例的 YAML 均放在该文件的同一目录。相对日志路径从启动工作目录计算。`go run . -h` 显示用法。
