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

func (s *Socket) Run() (err error) {
	err = s.conn.Connect()
	if err != nil {
		return
	}
	go s.conn.Run()
	go s.listen()
	return
}

func (s *Socket) Wright(w http.ResponseWriter, r *http.Request) {
	msg := Payload{}
	json.NewDecoder(r.Body).Decode(&msg)
	err := s.conn.SendToDevice(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.conn.Error <- err.Error()
	}
}

func (s *Socket) Read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(s.conn.ReadFromDevice())
}

func (s *Socket) listen() {
	for {
		if broke := <-s.conn.BrokenPipe; broke {
			select {
			case <-time.After(time.Second):
				log.Println("Broken pipe")
				s.try()
				return
			}
		}
	}
}

func (s *Socket) try() {
	err := s.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
