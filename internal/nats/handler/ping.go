package handler

import (
	"context"

	"github.com/nats-io/nats.go"
)

// HandlePing 向请求自带的回复主题发送 pong。
func HandlePing(ctx context.Context, msg *nats.Msg) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return msg.Respond([]byte("pong"))
}
