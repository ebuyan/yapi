package server

import (
	"encoding/json"
	"net/http"
	"yapi/socket"
)

type Handler struct {
	*socket.Conversation
}

func NewHandler(conversation *socket.Conversation) Handler {
	return Handler{conversation}
}

func (handler Handler) GetLastState(w http.ResponseWriter, r *http.Request) {
	js, _ := json.Marshal(handler.Conversation.Device.LastState)
	w.Header().Add("Content-Type", "application/json")
	w.Write(js)
}

func (handler Handler) SendCommand(w http.ResponseWriter, r *http.Request) {
	var msg map[string]interface{}
	json.NewDecoder(r.Body).Decode(&msg)
	handler.Conversation.SendToStation(msg)
}
