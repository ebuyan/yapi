package socket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Socket struct {
	conn *Conversation
}

func NewSocket(conn *Conversation) *Socket {
	return &Socket{conn}
}

func (s *Socket) Run() (err error) {
	log.Println("Start socket connection")

	ctx := context.Background()
	if err = s.conn.Connect(ctx); err != nil {
		return errors.New("connect error: " + err.Error())
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err = s.conn.Run(ctx); err != nil {
			cancel()
			log.Fatalln(err)
		}
	}()
	return
}

func (s *Socket) Write(w http.ResponseWriter, r *http.Request) {
	msg := Payload{}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		return
	}

	if err := s.conn.SendToDevice(msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.conn.Error <- "Write error: " + err.Error()
	}
}

func (s *Socket) Read(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(s.conn.ReadFromDevice())
}
