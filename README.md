# Baize Demo

本阶段演示白泽的进程初始化。需要 Go 1.27.3 和同级目录 `../baize`。首次运行时，从模板创建本地配置。

```sh
cp -R conf.example conf
```

日常修改 `conf/app.yaml`，该目录由 Git 忽略。运行 demo 并查看初始化日志。

```sh
go run .
cat log/baize-demo/biz/info.log
```

白泽读取 `./conf/app.yaml`，初始化 Pixiu 环境、时区和日志。默认日志写入“白泽初始化完成”后，程序关闭日志并退出。

可以通过 `go run . --conf /etc/cards/staging.yaml` 指定入口（简写为 `-c`），配置目录取该文件的父目录。相对的日志路径以进程启动工作目录为基准。`go run . -h` 显示用法。
