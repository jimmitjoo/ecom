package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/events"
	"github.com/jimmitjoo/ecom/src/domain/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jimmitjoo/ecom/src/infrastructure/logging"
	"go.uber.org/zap"
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
	writeMu   sync.Mutex // New mutex for write operations
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

// writeMessage is a thread-safe wrapper for writing to a WebSocket connection
func (h *WebSocketHandler) writeMessage(conn *websocket.Conn, messageType int, data []byte) error {
	h.writeMu.Lock()
	defer h.writeMu.Unlock()
	return conn.WriteMessage(messageType, data)
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	logger, _ := logging.NewLogger()
	logger = logger.WithRequestID(requestID)

	logger.Debug("New WebSocket connection attempt",
		zap.String("remote_addr", r.RemoteAddr),
	)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("WebSocket upgrade failed",
			zap.Error(err),
			zap.String("remote_addr", r.RemoteAddr),
		)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	clientCount := len(h.clients)
	h.mu.Unlock()

	logger.Info("New WebSocket client connected",
		zap.String("remote_addr", r.RemoteAddr),
		zap.Int("total_clients", clientCount),
	)

	// Clean up client when connection closes
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		clientCount := len(h.clients)
		h.mu.Unlock()
		conn.Close()

		logger.Info("WebSocket client disconnected",
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("remaining_clients", clientCount),
		)
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
			if err := h.writeMessage(conn, websocket.PongMessage, nil); err != nil {
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
	logger, _ := logging.NewLogger()

	startTime := time.Now()
	data, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal event for broadcast",
			zap.Error(err),
			zap.String("event_type", string(event.Type)),
			zap.String("event_id", event.ID),
		)
		return
	}

	h.mu.RLock()
	clientCount := len(h.clients)
	h.mu.RUnlock()

	logger.Debug("Starting event broadcast",
		zap.String("event_type", string(event.Type)),
		zap.String("event_id", event.ID),
		zap.Int("client_count", clientCount),
	)

	successCount := 0
	failCount := 0

	h.mu.RLock()
	clients := make([]*websocket.Conn, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		if err := h.writeMessage(client, websocket.TextMessage, data); err != nil {
			log.Printf("Failed to send message to client: %v", err)
			h.mu.Lock()
			delete(h.clients, client)
			h.mu.Unlock()
			client.Close()
			failCount++
		} else {
			successCount++
		}
	}

	logger.Info("Event broadcast completed",
		zap.String("event_type", string(event.Type)),
		zap.String("event_id", event.ID),
		zap.Int("success_count", successCount),
		zap.Int("fail_count", failCount),
		zap.Duration("duration", time.Since(startTime)),
	)
}
