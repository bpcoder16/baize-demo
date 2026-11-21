// Package handler 处理 Demo 的 Realtime 消息。
package handler

import (
	"context"

	"github.com/bpcoder16/baize/service/realtime"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Ping 接收空消息，用同一 requestId 回复 demo.pong。
func Ping(ctx context.Context, session *realtime.Session, msg *realtime.Message) error {
	if err := msg.Decode(&emptypb.Empty{}); err != nil {
		return session.ReplyError(ctx, msg, "demo.pong", realtime.CodeInvalidArgument, "ping 需要空消息")
	}
	return session.Reply(ctx, msg, "demo.pong", wrapperspb.String("pong"))
}
