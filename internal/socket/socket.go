package socket

import (
	"encoding/json"
	"log"
	"net/http"
)

type Socket struct {
	conn *Conversation
}

func NewSocket(conn *Conversation) Socket {
	return Socket{conn}
}

func (s *Socket) Run() (err error) {
	log.Println("Run socket")
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
	json.Marshal(msg)

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
		select {
		case <-s.conn.BrokenPipe:
			log.Println("Broken pipe")
			err := s.Run()
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
	}
}
