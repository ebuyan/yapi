package server

import (
	"net/http"
	"os"
	"time"
	"yapi/internal/socket"

	"github.com/gorilla/mux"
)

type Http struct {
	host   string
	socket *socket.Socket
}

func NewHttp(s *socket.Socket) Http {
	return Http{os.Getenv("HTTP_HOST"), s}
}

func (h *Http) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/", h.socket.Write).Methods("POST")
	r.HandleFunc("/", h.socket.Read).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         h.host,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	return srv.ListenAndServe()
}
