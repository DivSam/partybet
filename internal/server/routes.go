package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"partybet/internal/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.HelloWorldHandler)

	r.HandleFunc("/events/{id}", s.HandleNewEvent).Methods("POST")

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) HandleNewEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var event models.Event

	fmt.Printf("Request: %v\n", r.Body)

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("error upgrading connection. Err: %v", err)
	}

	defer conn.Close()

	eventJson, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, eventJson); err != nil {
		log.Fatalf("error writing message. Err: %v", err)
	}

}
