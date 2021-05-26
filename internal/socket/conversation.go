package socket

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Conversation struct {
	device     Device
	connection *websocket.Conn
	Error      chan string
	BrokenPipe chan bool
	writeWait  time.Duration
	pingPeriod time.Duration
}

func NewConversation(device Device) *Conversation {
	return &Conversation{
		device:     device,
		Error:      make(chan string),
		BrokenPipe: make(chan bool),
		writeWait:  10 * time.Second,
		pingPeriod: 300 * time.Second,
	}
}

func (c *Conversation) Connect() (err error) {
	dialer := websocket.DefaultDialer
	certs, err := GetCerts(c.device.GetCertificate())
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
	err = c.pingDevice()
	if err == nil {
		log.Println("Successful connection to the station")
	}
	return
}

func (c *Conversation) Run() {
	defer c.Close()
	go c.read()
	go c.pingConn()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case err := <-c.Error:
		log.Println(err)
		if strings.Contains(err, "Invalid token") {
			e := c.device.RefreshToken()
			if e != nil {
				log.Fatalln(e)
			}
		}
		select {
		case <-time.After(time.Second):
			for i := 0; i < 3; i++ {
				c.BrokenPipe <- true
			}
		}
	case <-interrupt:
		log.Fatalln("interrupt")
	}
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
	c.connection.SetWriteDeadline(time.Now().Add(c.writeWait))
	return c.connection.WriteJSON(message)
}

func (c *Conversation) Close() {
	c.connection.Close()
	log.Println("Connection closed")
}

func (c *Conversation) read() {
	log.Println("Start read socket")
	for {
		select {
		case <-c.BrokenPipe:
			return
		default:
			_, msg, err := c.connection.ReadMessage()
			if err != nil {
				c.Error <- "Read error: " + err.Error()
				return
			}
			c.device.SetState(msg)
		}
	}
}

func (c *Conversation) pingConn() {
	ticker := time.NewTicker(c.pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.connection.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(c.writeWait)); err != nil {
				c.Error <- "Ping error: " + err.Error()
				return
			}
		case <-c.BrokenPipe:
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
