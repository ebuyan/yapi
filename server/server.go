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
	Handler
}

func NewHttp(s *socket.Socket) Http {
	return Http{
		Host:    os.Getenv("HTTP_HOST"),
		Handler: NewHandler(s),
	}
}

func (h *Http) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", h.Handler.SetState).Methods("POST")
	r.HandleFunc("/", h.Handler.GetState).Methods("GET")
	http.Handle("/", r)
	log.Println("Start server on " + h.Host)
	log.Fatalln(http.ListenAndServe(h.Host, nil))
}
