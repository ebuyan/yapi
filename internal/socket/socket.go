package socket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type socket struct {
	conn *Conversation
}

func NewSocket(conn *Conversation) Socket {
	return &socket{conn}
}

func (s *socket) Run(ctx context.Context) (err error) {
	log.Println("start socket connection")

	if s.conn == nil {
		return errors.New("no socket connection")
	}

	if err = s.conn.Connect(ctx); err != nil {
		return errors.New("connect error: " + err.Error())
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err = s.conn.Run(ctx); err != nil {
			cancel()
			s.conn.Close()
			time.Sleep(time.Second * 1)
			log.Fatalln(err)
		}
	}()
	return
}

func (s *socket) Write(w http.ResponseWriter, r *http.Request) {
	msg := Payload{}
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		return
	}

	if err := s.conn.SendToDevice(msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *socket) Read(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(s.conn.ReadFromDevice())
}
