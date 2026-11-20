# Baize Demo

运行 `go run .` 会初始化白泽的环境、时区和日志，再依次连接 MySQL、Redis 和 NATS。三类资源的实例名均为 `default`。需要 Go 1.27.3、同级目录 `../baize` 和可连接的 MySQL、单节点 Redis、NATS 服务。

根入口通过 `mysqlboot.Instances("default")`、`redisboot.Instances("default")` 和 `natsboot.Instances("default")` 声明资源，并调用 `bootstrap.RunCLI(plan)`。首次运行或升级时，从模板补齐本地配置。

```sh
mkdir -p conf
cp -n conf.example/*.yaml conf/
```

应用配置在 `conf/app.yaml`。共享资源各有一份连接配置。

- `conf/mysql.default.yaml` 填写数据库地址、库名和凭据。
- `conf/redis.default.yaml` 填写主机、端口、认证信息和数据库编号。
- `conf/nats.default.yaml` 填写集群节点地址，启用认证时填写用户名和密码。根入口只建立客户端连接，可直接使用 Core NATS 服务完成启动验证。

三类资源就绪后，程序写入“白泽初始化完成”，随后排空并关闭 NATS，再关闭 Redis、MySQL 和日志。

```sh
go run .
cat log/baize-demo/biz/info.log
```

可用 `go run . --conf /etc/cards/staging.yaml` 指定应用配置（简写为 `-c`），三个实例的 YAML 均放在该文件的同一目录。相对日志路径从启动工作目录计算。`go run . -h` 显示用法。
