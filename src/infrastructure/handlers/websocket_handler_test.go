package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventPublisher struct {
	mock.Mock
	handlers map[models.EventType]func(*models.Event)
	mu       sync.RWMutex
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		handlers: make(map[models.EventType]func(*models.Event)),
	}
}

func (m *MockEventPublisher) Publish(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventPublisher) Subscribe(eventType models.EventType, handler func(*models.Event)) error {
	m.Called(eventType, handler)
	m.mu.Lock()
	m.handlers[eventType] = handler
	m.mu.Unlock()
	return nil
}

func (m *MockEventPublisher) Unsubscribe(eventType models.EventType, handler func(*models.Event)) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventPublisher) triggerHandler(eventType models.EventType, event *models.Event) {
	m.mu.RLock()
	handler, exists := m.handlers[eventType]
	m.mu.RUnlock()

	if exists {
		handler(event)
	}
}

func setupWebSocketTest() (*WebSocketHandler, *MockEventPublisher) {
	mockPublisher := NewMockEventPublisher()

	eventTypes := []models.EventType{
		models.EventProductCreated,
		models.EventProductUpdated,
		models.EventProductDeleted,
	}

	for _, eventType := range eventTypes {
		mockPublisher.On("Subscribe", eventType, mock.AnythingOfType("func(*models.Event)")).Return(nil)
	}

	return NewWebSocketHandler(mockPublisher), mockPublisher
}

func TestWebSocketConnection(t *testing.T) {
	handler, mockPublisher := setupWebSocketTest()

	server := httptest.NewServer(http.HandlerFunc(handler.HandleWebSocket))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Verifiera bara att anslutningen lyckades
	assert.NotNil(t, ws)

	mockPublisher.AssertExpectations(t)
}

func TestWebSocketBroadcast(t *testing.T) {
	handler, mockPublisher := setupWebSocketTest()

	server := httptest.NewServer(http.HandlerFunc(handler.HandleWebSocket))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	ws1, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws1.Close()

	ws2, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws2.Close()

	time.Sleep(100 * time.Millisecond)

	event := &models.Event{
		ID:   "test_event_1",
		Type: models.EventProductCreated,
		Data: &models.ProductEvent{
			ProductID: "test_prod_1",
			Action:    "created",
		},
	}

	mockPublisher.triggerHandler(models.EventProductCreated, event)

	for _, ws := range []*websocket.Conn{ws1, ws2} {
		done := make(chan *models.Event)
		go func(conn *websocket.Conn) {
			_, message, err := conn.ReadMessage()
			if err != nil {
				t.Errorf("Failed to read message: %v", err)
				return
			}
			var receivedEvent models.Event
			if err := json.Unmarshal(message, &receivedEvent); err != nil {
				t.Errorf("Failed to unmarshal event: %v", err)
				return
			}
			done <- &receivedEvent
		}(ws)

		select {
		case receivedEvent := <-done:
			assert.Equal(t, event.ID, receivedEvent.ID)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout waiting for event")
		}
	}

	mockPublisher.AssertExpectations(t)
}

func TestWebSocketClientDisconnect(t *testing.T) {
	handler, mockPublisher := setupWebSocketTest()

	server := httptest.NewServer(http.HandlerFunc(handler.HandleWebSocket))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)

	ws.Close()

	time.Sleep(100 * time.Millisecond)

	handler.mu.RLock()
	numClients := len(handler.clients)
	handler.mu.RUnlock()

	assert.Equal(t, 0, numClients)
	mockPublisher.AssertExpectations(t)
}
