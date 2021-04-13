package socket

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
	"yapi/glagol"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Conversation struct {
	Device     *glagol.Device
	Connection *websocket.Conn
	Error      chan string
	Interrupt  chan os.Signal
	Locked     bool
}

func NewConversation(device *glagol.Device) Conversation {
	return Conversation{
		Device:    device,
		Error:     make(chan string),
		Interrupt: make(chan os.Signal, 1),
		Locked:    false,
	}
}

func (conversation *Conversation) Run() {
	conversation.connect()
	go conversation.read()
	conversation.ping()
	conversation.finish()
}

func (conversation *Conversation) connect() {
	host := url.URL{Scheme: "wss", Host: conversation.Device.Discovery.GetHost(), Path: "/"}
	dialer := websocket.DefaultDialer

	certs, err := GetCerts(conversation.Device.Glagol.Security.ServerCertificate)
	if err != nil {
		panic(err)
	}
	dialer.TLSClientConfig = &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}
	dialer.HandshakeTimeout = 0
	conversation.Connection, _, err = dialer.Dial(host.String(), http.Header{"Origin": {"http://yandex.ru/"}})
	if err != nil {
		panic("dial: " + err.Error())
	}
}

func (conversation *Conversation) read() {
	for {
		_, message, err := conversation.Connection.ReadMessage()
		if err != nil {
			conversation.Error <- err.Error()
			return
		}
		go conversation.updateState(message)
	}
}

func (conversation *Conversation) updateState(msg []byte) {
	for conversation.Locked {
	}
	conversation.Locked = true

	latestState := glagol.DeviceResponse{}
	json.Unmarshal(msg, &latestState)
	conversation.Device.LastState = latestState

	conversation.Locked = false
}

func (conversation *Conversation) ping() {
	err := conversation.SendToStation(map[string]interface{}{"command": "ping"})
	if err != nil {
		conversation.Error <- "write: " + err.Error()
	}
}

func (conversation *Conversation) finish() {
	defer conversation.Connection.Close()
	signal.Notify(conversation.Interrupt, os.Interrupt)
	for {
		select {
		case err := <-conversation.Error:
			panic(err)
		case <-conversation.Interrupt:
			return
		}
	}
}

func (conversation *Conversation) SendToStation(msg interface{}) error {
	payload := DeviceRequest{
		ConversationToken: conversation.Device.Token,
		Id:                uuid.New().String(),
		SentTime:          time.Now().UnixNano(),
		Payload:           msg,
	}
	return conversation.Connection.WriteJSON(payload)
}

type DeviceRequest struct {
	ConversationToken string      `json:"conversationToken"`
	Id                string      `json:"id"`
	SentTime          int64       `json:"sentTime"`
	Payload           interface{} `json:"payload"`
}

type Payload struct {
	Command string `json:"command"`
	Volume  string `json:"volume"`
	Text    string `json:"text"`
}
