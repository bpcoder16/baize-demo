# Baize Demo

运行 `go run .` 会初始化白泽的环境、时区和日志，再依次连接 MySQL 和 Redis。两类资源的实例名均为 `default`。需要 Go 1.27.3、同级目录 `../baize` 和可连接的 MySQL、单节点 Redis。

根入口通过 `mysqlboot.Instances("default")` 和 `redisboot.Instances("default")` 声明资源，并调用 `bootstrap.RunCLI(plan)`。首次运行或升级时，从模板补齐本地配置。

```sh
mkdir -p conf
cp -n conf.example/*.yaml conf/
```

在 `conf/mysql.default.yaml` 填写数据库地址、库名和凭据，在 `conf/redis.default.yaml` 填写 Redis 的主机、端口、认证信息和数据库编号。应用配置在 `conf/app.yaml`。Redis 通过 `PING` 验活后，程序写入“白泽初始化完成”，依次关闭 Redis、MySQL 和日志，然后退出。

```sh
go run .
cat log/baize-demo/biz/info.log
```

可用 `go run . --conf /etc/cards/staging.yaml` 指定应用配置（简写为 `-c`），`mysql.default.yaml` 和 `redis.default.yaml` 均放在该文件的同一目录。相对日志路径从启动工作目录计算。`go run . -h` 显示用法。
