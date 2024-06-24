package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"partybet/internal/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[int]map[*websocket.Conn]bool)
	broadcast = make(chan models.Bet)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var (
	events      = make(map[int]*models.Event)
	eventsMutex sync.Mutex
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.HelloWorldHandler)

	r.HandleFunc("/event", s.HandleNewEvent).Methods("POST")
	r.HandleFunc("/ws", s.HandleNewWebsocketConnection)

	go HandleBroadcast()

	fmt.Println("Server is running on port: ", s.port)

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
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: have some error handling for the payload, it needs to have all the fields (maybe we handle this in the client application)

	eventsMutex.Lock()
	events[event.ID] = &event // will use UUID or something in the client application
	eventsMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)

}

func (s *Server) HandleNewWebsocketConnection(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatalf("error upgrading connection to websocket. Err: %v", err)
	}

	defer ws.Close()

	eventId, err := strconv.Atoi(r.URL.Query().Get("event_id"))
	if err != nil {
		log.Fatalf("error parsing event_id. Err: %v", err)
	}

	if clients[eventId] == nil {
		clients[eventId] = make(map[*websocket.Conn]bool)
	}
	clients[eventId][ws] = true

	for {
		var bet models.Bet
		err := ws.ReadJSON(&bet)

		fmt.Println("bet: ", bet)

		if err != nil {
			log.Printf("error reading JSON. Err: %v", err)
			break
		}

		// fmt.Println("clients: ", clients)
		// we send the bet to the broadcasts
		broadcast <- bet

	}
}

func HandleBroadcast() {
	for {
		bet := <-broadcast

		event := events[bet.EventID]

		eventsMutex.Lock()
		//TODO: handle core betting logic here
		event.UpdateHandle(bet.Outcome, bet.Amount)

		eventsMutex.Unlock()

		fmt.Println("event: ", event.Total, event.HandleYes, event.HandleNo)
		fmt.Println(len(clients[bet.EventID]))
		// send event to all clients with the same event id
		for client := range clients[bet.EventID] {
			err := client.WriteJSON(event)
			if err != nil {
				log.Printf("error writing JSON. Err: %v", err)
				client.Close()
				delete(clients[bet.EventID], client)
			}
		}

	}
}
