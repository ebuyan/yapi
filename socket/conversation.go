package socket

import (
	"crypto/tls"
	"encoding/json"
	"log"
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
	BrokenPipe chan bool
	Locked     bool
}

func NewConversation(device *glagol.Device) *Conversation {
	return &Conversation{
		Device:     device,
		Error:      make(chan string),
		BrokenPipe: make(chan bool),
		Locked:     false,
	}
}

func (c *Conversation) Connect() (err error) {
	dialer := websocket.DefaultDialer
	certs, err := GetCerts(c.Device.Glagol.Security.ServerCertificate)
	if err != nil {
		return
	}
	dialer.TLSClientConfig = &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}
	c.Connection, _, err = dialer.Dial(c.Device.Config.GetHost(), c.Device.Config.GetHeaderOrigin())
	if err != nil {
		return
	}
	log.Println("Successful connection to the station")
	err = c.ping()
	return
}

func (c *Conversation) Run() {
	c.BrokenPipe <- false
	defer c.Close()
	go c.read()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case err := <-c.Error:
		log.Println(err)
		select {
		case <-time.After(time.Second):
			c.BrokenPipe <- true
		}
	case <-interrupt:
		log.Fatalln("interrupt")
	}
}

func (c *Conversation) SendToDevice(msg Payload) error {
	payload := DeviceRequest{
		ConversationToken: c.Device.Token,
		Id:                uuid.New().String(),
		SentTime:          time.Now().UnixNano(),
		Payload:           msg,
	}
	return c.Connection.WriteJSON(payload)
}

func (c *Conversation) Close() {
	c.Connection.Close()
	log.Println("Connection closed")
}

func (c *Conversation) read() {
	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			c.Error <- err.Error()
			return
		}
		go c.updateState(message)
	}
}

func (c *Conversation) updateState(msg []byte) {
	for c.Locked {
	}
	c.Locked = true

	latestState := glagol.DeviceState{}
	json.Unmarshal(msg, &latestState)
	c.Device.State = latestState

	c.Locked = false
}

func (c *Conversation) ping() (err error) {
	err = c.SendToDevice(Payload{Command: "ping"})
	if err != nil {
		return
	}
	return
}

type DeviceRequest struct {
	ConversationToken string  `json:"conversationToken"`
	Id                string  `json:"id"`
	SentTime          int64   `json:"sentTime"`
	Payload           Payload `json:"payload"`
}

type Payload struct {
	Command  string  `json:"command"`
	Volume   float32 `json:"volume"`
	Position int8    `json:"position"`
	Text     string  `json:"text"`
}
