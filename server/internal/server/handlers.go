package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"partybet/server/internal/models"

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

		// the hack here is going to be to send an empty ws message to the server upon creation if the creator doesn't want to bet
		if event, exists := events[eventId]; exists {
			s.StartEventTimer(eventId, event.Duration)
		}
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

func (s *Server) StartEventTimer(eventId int, duration string) {
	realDuration, _ := time.ParseDuration(duration)
	time.AfterFunc(realDuration, func() {
		if conns, ok := clients[eventId]; ok {
			for ws := range conns {
				ws.Close()
			}
			fmt.Println("event finished, closing all connections")
			delete(clients, eventId)
		}
	})
}
