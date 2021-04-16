package socket

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Conversation struct {
	device     Device
	connection *websocket.Conn
	Error      chan string
	BrokenPipe chan bool
}

func NewConversation(device Device) *Conversation {
	return &Conversation{
		device:     device,
		Error:      make(chan string),
		BrokenPipe: make(chan bool),
	}
}

func (c *Conversation) Connect() (err error) {
	dialer := websocket.DefaultDialer
	certs, err := GetCerts(c.device.GetSertificate())
	if err != nil {
		return
	}
	dialer.TLSClientConfig = &tls.Config{
		RootCAs:            certs,
		InsecureSkipVerify: true,
	}
	c.connection, _, err = dialer.Dial(c.device.GetHost(), c.device.GetOrigin())
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

func (c *Conversation) ReadFromDevice() []byte {
	for c.device.Locked() {
	}
	return c.device.GetState()
}

func (c *Conversation) SendToDevice(msg Payload) error {
	message := DeviceRequest{
		ConversationToken: c.device.GetToken(),
		Id:                uuid.New().String(),
		SentTime:          time.Now().UnixNano(),
		Payload:           msg,
	}
	return c.connection.WriteJSON(message)
}

func (c *Conversation) Close() {
	c.connection.Close()
	log.Println("Connection closed")
}

func (c *Conversation) read() {
	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			c.Error <- err.Error()
			return
		}
		go c.updateState(message)
	}
}

func (c *Conversation) updateState(msg []byte) {
	for c.device.Locked() {
	}
	c.device.Lock()
	defer c.device.Unlock()
	c.device.SetState(msg)
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
