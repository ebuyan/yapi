package socket

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type device struct {
	host string
}

func (d device) GetState() []byte                   { return nil }
func (d device) SetState(_ []byte)                  {}
func (d device) GetHost() string                    { return d.host }
func (d device) RefreshToken(context.Context) error { return nil }
func (d device) GetToken() string                   { return "" }
func (d device) GetCertificate() string             { return testCert }

func TestConversation_Connect(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	conn := NewConversation(device{makeWsProto(s.URL)})

	err := conn.Connect(context.Background())
	defer conn.Close()

	require.Nil(t, err)
}

func TestConversation_Run(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	conn := NewConversation(device{makeWsProto(s.URL)})

	err := conn.Connect(context.Background())
	defer conn.Close()

	require.Nil(t, err)

	go func() {
		err = conn.Run(context.Background())
		require.NotNil(t, err)
	}()

	conn.error <- "err"
}

func newServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(cstHandler{
		t:        t,
		upgrader: websocket.Upgrader{},
	})
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

type cstHandler struct {
	t        *testing.T
	upgrader websocket.Upgrader
}

func (h cstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(w, r, nil)
	require.Nil(h.t, err)
	defer func() { _ = ws.Close() }()
}
