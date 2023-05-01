package server

import (
	"net/http"
	"time"
	"yapi/internal/socket"

	"github.com/gorilla/mux"
)

type Http struct {
	addr   string
	socket socket.Socket
}

func NewHttp(s socket.Socket, addr string) Http {
	return Http{addr, s}
}

func (h *Http) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/", h.socket.Write).Methods("POST")
	r.HandleFunc("/", h.socket.Read).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         h.addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	return srv.ListenAndServe()
}
