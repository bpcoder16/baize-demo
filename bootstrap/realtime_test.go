package bootstrap_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/bpcoder16/baize-demo/bootstrap"
	realtimeinfra "github.com/bpcoder16/baize-demo/internal/infra/realtime"
	baize "github.com/bpcoder16/baize/bootstrap"
	"github.com/bpcoder16/baize/service/realtime"
	"github.com/bpcoder16/baize/service/realtime/protocol"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestDemoRealtime(t *testing.T) {
	// 命名实例与宿主进程同寿命，每次测试用独立的 app 实例。
	if os.Getenv("BAIZE_DEMO_TEST_REALTIME") == "" {
		ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDemoRealtime$")
		cmd.Env = append(os.Environ(), "BAIZE_DEMO_TEST_REALTIME=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%v\n%s", err, output)
		}
		return
	}
	services := bootstrap.NewPlan().Services
	i := slices.IndexFunc(services, func(s baize.Service) bool { return s.Name == "realtime:app" })
	if i < 0 {
		t.Fatal("启动计划缺少 realtime:app")
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	template, err := os.ReadFile("../conf.example/realtime.app.yaml")
	if err != nil {
		t.Fatal(err)
	}
	want := []byte("127.0.0.1:8889")
	if !bytes.Contains(template, want) {
		t.Fatal("模板缺少本机监听地址")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "realtime.app.yaml"), bytes.Replace(template, want, []byte(addr), 1), 0600); err != nil {
		t.Fatal(err)
	}
	s, err := services[i].Build(t.Context(), dir)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- s.Run() }()
	t.Cleanup(func() {
		s.Cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Error(err)
			}
		case <-time.After(3 * time.Second):
			t.Error("Run 未退出")
		}
	})
	client := http.Client{Timeout: time.Second}
	base := "http://" + addr
	deadline := time.Now().Add(3 * time.Second)
	for {
		resp, err := client.Get(base + "/ready/ws")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusNoContent {
				break
			}
		}
		if time.Now().After(deadline) {
			t.Fatal("就绪探针超时")
		}
		time.Sleep(5 * time.Millisecond)
	}
	for path, status := range map[string]int{"/ready/ws": 405, "/ws/app": 405, "/app/ping": 404} {
		resp, err := client.Head(base + path)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != status {
			t.Fatalf("HEAD %s = %d", path, resp.StatusCode)
		}
	}
	for _, sub := range []string{realtime.SubprotocolJSON, realtime.SubprotocolProtobuf} {
		dialer := websocket.Dialer{Subprotocols: []string{sub}, HandshakeTimeout: time.Second}
		conn, _, err := dialer.Dial("ws://"+addr+"/ws/app", nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			t.Fatal(err)
		}
		if conn.Subprotocol() != sub {
			t.Fatalf("握手子协议 %s", conn.Subprotocol())
		}
		if sub == realtime.SubprotocolJSON {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"kind":1,"requestId":"r","scene":"demo","action":"ping","data":{}}`)); err != nil {
				t.Fatal(err)
			}
			frame, wire, err := conn.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}
			var msg struct {
				Kind                  int
				RequestID, Type, Data string
			}
			if err := json.Unmarshal(wire, &msg); err != nil {
				t.Fatal(err)
			}
			if frame != websocket.TextMessage || msg.Kind != 2 || msg.RequestID != "r" || msg.Type != "demo.pong" || msg.Data != "pong" {
				t.Fatalf("意外响应 %s", wire)
			}
		} else {
			req := &protocol.Envelope{
				Kind:      protocol.Kind_REQUEST,
				RequestId: "r",
				Scene:     "demo",
				Action:    "ping",
				Body:      &protocol.Envelope_Data{},
			}
			wire, err := proto.Marshal(req)
			if err != nil {
				t.Fatal(err)
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, wire); err != nil {
				t.Fatal(err)
			}
			frame, wire, err := conn.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}
			var msg protocol.Envelope
			if err := proto.Unmarshal(wire, &msg); err != nil {
				t.Fatal(err)
			}
			var pong wrapperspb.StringValue
			if err := proto.Unmarshal(msg.GetData(), &pong); err != nil {
				t.Fatal(err)
			}
			if frame != websocket.BinaryMessage || msg.Kind != protocol.Kind_RESPONSE || msg.RequestId != "r" || msg.Type != "demo.pong" || pong.Value != "pong" {
				t.Fatalf("意外响应 %v", &msg)
			}
		}
	}
	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	if err := s.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	if realtimeinfra.App().Ready() {
		t.Fatal("停收后仍就绪")
	}
}
