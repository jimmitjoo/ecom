package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/jimmitjoo/ecom/src/domain/events"
	"github.com/jimmitjoo/ecom/src/domain/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, this should be restricted
	},
}

type WebSocketHandler struct {
	clients   map[*websocket.Conn]bool
	publisher events.EventPublisher
	mu        sync.RWMutex
}

func NewWebSocketHandler(publisher events.EventPublisher) *WebSocketHandler {
	handler := &WebSocketHandler{
		clients:   make(map[*websocket.Conn]bool),
		publisher: publisher,
	}

	// Subscribe to all product events
	handler.subscribeToEvents()

	return handler
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Websocket upgrade failed: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	// Clean up client when connection closes
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
	}()

	// Keep connection open and handle messages
	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket error: %v", err)
			}
			break
		}

		if messageType == websocket.PingMessage {
			if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				log.Printf("Failed to send pong: %v", err)
				break
			}
		}
	}
}

func (h *WebSocketHandler) subscribeToEvents() {
	eventTypes := []models.EventType{
		models.EventProductCreated,
		models.EventProductUpdated,
		models.EventProductDeleted,
	}

	for _, eventType := range eventTypes {
		h.publisher.Subscribe(eventType, func(event *models.Event) {
			h.broadcastEvent(event)
		})
	}
}

func (h *WebSocketHandler) broadcastEvent(event *models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if err := client.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Failed to send message to client: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
