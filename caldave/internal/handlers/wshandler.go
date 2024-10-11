package handlers

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Client struct {
	connection *websocket.Conn
	hub        *Hub
}

type Hub struct {
	clientsList map[*websocket.Conn]bool
	broadcast   chan []byte
	mutex       sync.RWMutex
}

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler() *WebSocketHandler {
	hub := &Hub{
		clientsList: make(map[*websocket.Conn]bool),
		broadcast:   make(chan []byte),
	}
	go hub.run()
	return &WebSocketHandler{hub: hub}
}

func (h *Hub) run() {
	for message := range h.broadcast {
		h.mutex.RLock()
		for conn := range h.clientsList {
			err := websocket.Message.Send(conn, string(message))
			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
				h.mutex.RUnlock()
				h.removeClient(conn)
				h.mutex.RLock()
			}
		}
		h.mutex.RUnlock()
	}
}

func (h *Hub) addClient(conn *websocket.Conn) {
	h.mutex.Lock()
	h.clientsList[conn] = true
	h.mutex.Unlock()
}

func (h *Hub) removeClient(conn *websocket.Conn) {
	h.mutex.Lock()
	delete(h.clientsList, conn)
	conn.Close()
	h.mutex.Unlock()

	// h.mutex.Lock()
	// defer h.mutex.Unlock()
	// if _, ok := h.clientsList[conn]; ok {
	// 	conn.Close()
	// 	delete(h.clientsList, conn)
	// }
}

func (wsh *WebSocketHandler) handleWS(ws *websocket.Conn) {
	log.Printf("New connection from %s", ws.RemoteAddr())
	wsh.hub.addClient(ws)
	defer wsh.hub.removeClient(ws)

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
