package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"yapi/glagol"
)

type Socket struct {
	*Conversation
}

func NewSocket(device *glagol.Device) Socket {
	return Socket{NewConversation(device)}
}

func (s Socket) Run() (err error) {
	err = s.Conversation.Connect()
	if err != nil {
		return
	}
	go s.listen()
	return
}

func (s Socket) Wright(w http.ResponseWriter, r *http.Request) {
	msg := Payload{}
	json.NewDecoder(r.Body).Decode(&msg)
	err := s.Conversation.SendToDevice(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.Conversation.Error <- err.Error()
	}
}

func (s Socket) Read(w http.ResponseWriter) {
	js, _ := json.Marshal(s.Conversation.Device.State)
	w.Write(js)
}

func (s Socket) listen() {
	go s.Conversation.Run()

	for {
		broke := <-s.Conversation.BrokenPipe
		if broke {
			log.Println("Broken pipe")
			err := s.Run()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
