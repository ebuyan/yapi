package server

import (
	"net/http"
	"yapi/socket"
)

type Handler struct {
	*socket.Socket
}

func NewHandler(s *socket.Socket) Handler {
	return Handler{s}
}

func (h Handler) GetState(w http.ResponseWriter, r *http.Request) {
	h.Socket.Read(w)
	w.Header().Add("Content-Type", "application/json")
}

func (h Handler) SetState(w http.ResponseWriter, r *http.Request) {
	h.Socket.Wright(w, r)
}
