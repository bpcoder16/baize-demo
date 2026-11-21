package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bpcoder16/pixiu/logit"
	"github.com/nats-io/nats.go/jetstream"
)

// HandleDemoEvent 将事件 ID 加入本次消费日志，确认和重投由 natsx 处理。
func HandleDemoEvent(ctx context.Context, msg jetstream.Msg) error {
	var event struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return fmt.Errorf("decode demo event: %w", err)
	}
	if strings.TrimSpace(event.ID) == "" {
		return errors.New("demo event requires id")
	}
	logit.AddField(ctx, logit.Str("event_id", event.ID))
	return nil
}
