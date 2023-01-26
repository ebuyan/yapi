package socket

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait  = 60 * time.Second
	pingWait  = (pongWait * 9) / 10
)

type Conversation struct {
	device     Device
	connection *websocket.Conn
	Error      chan string
}

func NewConversation(device Device) *Conversation {
	return &Conversation{
		device: device,
		Error:  make(chan string),
	}
}

func (c *Conversation) Connect(ctx context.Context) (err error) {
	dialer := websocket.DefaultDialer
	certs, err := GetCerts(c.device.GetCertificate())
	if err != nil {
		return
	}

	dialer.TLSClientConfig = &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}

	if c.connection, _, err = dialer.DialContext(ctx, c.device.GetHost(), c.device.GetOrigin()); err != nil {
		return
	}
	if err = c.pingDevice(); err == nil {
		log.Println("successful connection to the station")
	}
	return
}

func (c *Conversation) Run(ctx context.Context) error {
	defer c.Close()
	go c.read(ctx)
	go c.pingConn(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case e := <-c.Error:
		log.Println(e)
		if strings.Contains(e, "Invalid token") {
			if err := c.device.RefreshToken(); err != nil {
				return err
			}
		}
	case <-interrupt:
		return errors.New("interrupt")
	}
	return nil
}

func (c *Conversation) ReadFromDevice() []byte {
	return c.device.GetState()
}

func (c *Conversation) SendToDevice(msg Payload) error {
	message := DeviceRequest{
		ConversationToken: c.device.GetToken(),
		Id:                uuid.New().String(),
		SentTime:          time.Now().UnixNano(),
		Payload:           msg,
	}
	_ = c.connection.SetWriteDeadline(time.Now().Add(writeWait))
	return c.connection.WriteJSON(message)
}

func (c *Conversation) Close() {
	_ = c.connection.Close()
	log.Println("connection closed")
}

func (c *Conversation) read(ctx context.Context) {
	log.Println("start read socket")
	c.connection.SetReadDeadline(time.Now().Add(pongWait))

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := c.connection.ReadMessage()
			if err != nil {
				c.Error <- "read error: " + err.Error()
				return
			}
			c.device.SetState(msg)
		}
	}
}

func (c *Conversation) pingConn(ctx context.Context) {
	ticker := time.NewTicker(pingWait)
	defer ticker.Stop()

	c.connection.SetPongHandler(func(string) error { c.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		select {
		case <-ticker.C:
			if err := c.pingDevice(); err != nil {
				c.Error <- "ping error: " + err.Error()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Conversation) pingDevice() (err error) {
	return c.SendToDevice(Payload{Command: "ping"})
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
