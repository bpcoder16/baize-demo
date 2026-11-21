package handler

import "context"

// Heartbeat 用于观察调度是否运行，每轮结果由 Cron 统一记录。
func Heartbeat(ctx context.Context) error {
	return ctx.Err()
}
