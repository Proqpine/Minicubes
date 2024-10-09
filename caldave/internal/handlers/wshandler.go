package handlers

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Hub struct {
	connections map[*websocket.Conn]bool
	broadcast   chan []byte
	mutex       sync.RWMutex
}

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler() *WebSocketHandler {
	hub := &Hub{
		connections: make(map[*websocket.Conn]bool),
		broadcast:   make(chan []byte),
	}
	go hub.run()
	return &WebSocketHandler{hub: hub}
}

func (h *Hub) run() {
	for message := range h.broadcast {
		h.mutex.RLock()
		for conn := range h.connections {
			err := websocket.Message.Send(conn, string(message))
			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
				h.mutex.RUnlock()
				h.removeConnection(conn)
				h.mutex.RLock()
			}
		}
		h.mutex.RUnlock()
	}
}

func (h *Hub) addConnection(conn *websocket.Conn) {
	h.mutex.Lock()
	h.connections[conn] = true
	h.mutex.Unlock()
}

func (h *Hub) removeConnection(conn *websocket.Conn) {
	h.mutex.Lock()
	delete(h.connections, conn)
	conn.Close()
	h.mutex.Unlock()
}

func (wsh *WebSocketHandler) handleWS(ws *websocket.Conn) {
	log.Printf("New connection from %s", ws.RemoteAddr())
	wsh.hub.addConnection(ws)
	defer wsh.hub.removeConnection(ws)

	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			return
		}

		log.Printf("Recieved message: %s", message)
		wsh.hub.broadcast <- []byte(message)
	}
}

func (wsh *WebSocketHandler) Handler() http.Handler {
	return websocket.Handler(wsh.handleWS)
}
