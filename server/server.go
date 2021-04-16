package server

import (
	"log"
	"net/http"
	"os"
	"yapi/socket"

	"github.com/gorilla/mux"
)

type Http struct {
	Host string
	*socket.Socket
}

func NewHttp(s *socket.Socket) Http {
	return Http{os.Getenv("HTTP_HOST"), s}
}

func (h *Http) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", h.Socket.Wright).Methods("POST")
	r.HandleFunc("/", h.Socket.Read).Methods("GET")
	http.Handle("/", r)
	log.Println("Start server on " + h.Host)
	log.Fatalln(http.ListenAndServe(h.Host, nil))
}
