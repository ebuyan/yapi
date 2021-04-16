package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Socket struct {
	conn *Conversation
}

func NewSocket(conn *Conversation) Socket {
	return Socket{conn}
}

func (s Socket) Run() (err error) {
	err = s.conn.Connect()
	if err != nil {
		return
	}
	go s.listen()
	return
}

func (s Socket) Wright(w http.ResponseWriter, r *http.Request) {
	msg := Payload{}
	json.NewDecoder(r.Body).Decode(&msg)
	err := s.conn.SendToDevice(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.conn.Error <- err.Error()
	}
}

func (s Socket) Read(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(s.conn.ReadFromDevice())
	w.Header().Add("Content-Type", "application/json")
	w.Write(js)
}

func (s Socket) listen() {
	go s.conn.Run()

	for {
		broke := <-s.conn.BrokenPipe
		if broke {
			log.Println("Broken pipe")
			s.waitDevice()
		}
	}
}

func (s Socket) waitDevice() {
	err := s.Run()
	if err != nil {
		log.Println("Wait device: " + err.Error())
		select {
		case <-time.After(time.Second):
			s.waitDevice()
		}
	}
}
