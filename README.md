# Baize Demo

运行 `go run .` 会初始化环境、时区和日志，打开 MySQL、PostgreSQL、ClickHouse、SQLite、Redis、NATS 和 Elasticsearch，再在本地构造阿里云 OSS 和 STS 客户端。需要 Go 1.27.3、同级目录 `../baize`、可用的 C 编译器，以及配置中指定的数据库和网络服务。SQLite 默认使用进程内存。

根入口在 `Plan.Resources` 中声明各实例，并调用 `bootstrap.RunCLI(plan)`。首次运行或升级时，从模板补齐本地配置。

```sh
mkdir -p conf
cp -n conf.example/*.yaml conf/
```

应用配置在 `conf/app.yaml`。共享资源各有一份连接配置。

- `conf/mysql.default.yaml` 填写数据库地址、库名和凭据。
- `conf/postgresql.default.yaml` 填写主库、可选从库及凭据。模板使用 `sslMode: "verify-full"`，本地配置按实际服务填写 TLS 模式和证书。
- `conf/clickhouse.default.yaml` 填写库名、凭据、协议和地址。Native 默认端口 9000，HTTP 默认端口 8123。
- `conf/sqlite.default.yaml` 默认使用 `:memory:`。持久化时将 `dsn` 改为文件路径或 `file:` URI，并提前创建父目录。构建 demo 时使用 `CGO_ENABLED=1`。
- `conf/redis.default.yaml` 填写主机、端口、认证信息和数据库编号。
- `conf/nats.default.yaml` 填写集群节点地址，启用认证时填写用户名和密码。根入口只建立客户端连接，可直接使用 Core NATS 服务完成启动验证。
- `conf/elasticsearch.default.yaml` 填写节点地址和认证信息。根入口导入 `elasticsearchboot/v8`，启动时通过 `GET /` 校验集群为 Elasticsearch 8。升级旧配置时删除 `version` 字段。使用自建 CA 时填写 `caCertFile`，相对路径从该文件所在目录解析。
- `conf/aliyunossx.default.yaml` 配置地域、Bucket 和访问密钥，业务通过 `ossx.Named("default")` 访问该桶。
- `conf/aliyunstsx.oss.yaml` 配置 STS 地域和访问密钥，业务通过 `stsx.Named("oss")` 申请临时凭据。

OSS 和 STS 的示例凭据可用于启动 demo。实际访问阿里云时，在本地配置中填写获授权的凭据。资源就绪后，程序写入“白泽初始化完成”，按声明的逆序回收客户端，最后关闭日志。所有资源共享 `conf/app.yaml` 的 `resourceCloseTimeout`（默认 30 秒），单个 NATS 客户端的 `closeTimeout`（默认 20 秒）计入这段时间。

```sh
go run .
cat log/baize-demo/biz/info.log
```

可用 `go run . --conf /etc/cards/staging.yaml` 指定应用配置（简写为 `-c`），各实例的 YAML 均放在该文件的同一目录。相对日志路径从启动工作目录计算。`go run . -h` 显示用法。
