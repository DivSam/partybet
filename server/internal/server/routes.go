package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/event", s.HandleNewEvent).Methods("POST")
	r.HandleFunc("/ws", s.HandleNewWebsocketConnection)

	go HandleBroadcast()

	fmt.Println("Server is running on port: ", s.port)

	return r
}
