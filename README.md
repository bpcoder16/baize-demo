# Baize Demo

运行 `go run .` 会初始化白泽的环境、时区和日志，再连接名为 `default` 的 MySQL 实例。需要 Go 1.27.3、同级目录 `../baize` 和可连接的 MySQL。

根入口在 `Plan.Resources` 声明 MySQL，并调用 `bootstrap.RunCLI(plan)`。首次运行时，从模板补齐本地配置。

```sh
mkdir -p conf
cp -n conf.example/*.yaml conf/
```

在 `conf/mysql.default.yaml` 填写数据库地址、库名和凭据，应用配置在 `conf/app.yaml`。运行成功后，程序写入“白泽初始化完成”，关闭连接池和日志并退出。

```sh
go run .
cat log/baize-demo/biz/info.log
```

可用 `go run . --conf /etc/cards/staging.yaml` 指定应用配置（简写为 `-c`），`mysql.default.yaml` 放在该文件的同一目录。相对日志路径从启动工作目录计算。`go run . -h` 显示用法。
