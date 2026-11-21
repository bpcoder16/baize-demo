// Package realtime 提供 Demo 的命名 Realtime 实例。
package realtime

import "github.com/bpcoder16/baize/service/realtime"

// AppName 对应 realtime.app.yaml。
const AppName = "app"

// App 供业务向本机 App 会话推送消息。
func App() *realtime.Server {
	return realtime.Named(AppName)
}
