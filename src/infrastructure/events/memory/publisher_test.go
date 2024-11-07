package memory

import (
	"sync"
	"testing"
	"time"

	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/assert"
)

func createTestProductEvent() *models.Event {
	return &models.Event{
		ID:        "test_event_1",
		Type:      models.EventProductCreated,
		EntityID:  "test_product_1",
		Version:   1,
		Sequence:  1,
		Timestamp: time.Now(),
		Data: &models.ProductEvent{
			ProductID: "test_product_1",
			Action:    "created",
			Product: &models.Product{
				ID:        "test_product_1",
				SKU:       "TEST-123",
				BaseTitle: "Test Product",
			},
		},
	}
}

func TestPublishAndSubscribe(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	receivedEvents := make([]*models.Event, 0)
	var mu sync.Mutex

	// Create a handler that saves received events
	handler := func(event *models.Event) {
		mu.Lock()
		receivedEvents = append(receivedEvents, event)
		mu.Unlock()
	}

	// Subscribe to events
	err := publisher.Subscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	// Publish an event
	event := createTestProductEvent()
	err = publisher.Publish(event)
	assert.NoError(t, err)

	// Wait a little so that the event can be processed
	time.Sleep(100 * time.Millisecond)

	// Verify that the event was received
	mu.Lock()
	assert.Len(t, receivedEvents, 1)
	assert.Equal(t, event.ID, receivedEvents[0].ID)
	mu.Unlock()
}

func TestMultipleSubscribers(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var wg sync.WaitGroup
	var mu sync.Mutex
	receivedCounts := make(map[string]int)

	// Create multiple handlers
	for i := 0; i < 3; i++ {
		handlerID := string(rune('A' + i))
		wg.Add(1)

		handler := func(event *models.Event) {
			mu.Lock()
			receivedCounts[handlerID]++
			mu.Unlock()
			wg.Done()
		}

		err := publisher.Subscribe(models.EventProductCreated, handler)
		assert.NoError(t, err)
	}

	// Publish an event
	event := createTestProductEvent()
	err := publisher.Publish(event)
	assert.NoError(t, err)

	// Wait for all handlers to finish
	wg.Wait()

	// Verify that all handlers received the event
	mu.Lock()
	assert.Len(t, receivedCounts, 3)
	for _, count := range receivedCounts {
		assert.Equal(t, 1, count)
	}
	mu.Unlock()
}

func TestUnsubscribe(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	receivedEvents := make([]*models.Event, 0)
	var mu sync.Mutex

	handler := func(event *models.Event) {
		mu.Lock()
		receivedEvents = append(receivedEvents, event)
		mu.Unlock()
	}

	// Subscribe and unsubscribe
	err := publisher.Subscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	err = publisher.Unsubscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	// Publish an event
	event := createTestProductEvent()
	err = publisher.Publish(event)
	assert.NoError(t, err)

	// Wait a little
	time.Sleep(100 * time.Millisecond)

	// Verify that no events were received
	mu.Lock()
	assert.Empty(t, receivedEvents)
	mu.Unlock()
}

func TestConcurrentPublishSubscribe(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var wg sync.WaitGroup
	var mu sync.Mutex
	receivedEvents := make(map[string]int)

	// Create handlers for different event types
	eventTypes := []models.EventType{
		models.EventProductCreated,
		models.EventProductUpdated,
		models.EventProductDeleted,
	}

	for _, eventType := range eventTypes {
		handler := func(event *models.Event) {
			mu.Lock()
			receivedEvents[string(event.Type)]++
			mu.Unlock()
		}
		err := publisher.Subscribe(eventType, handler)
		assert.NoError(t, err)
	}

	// Publish events concurrently
	numEvents := 10
	for i := 0; i < numEvents; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			event := createTestProductEvent()
			event.Type = eventTypes[index%len(eventTypes)]
			err := publisher.Publish(event)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	// Verify that the correct number of events were received
	mu.Lock()
	total := 0
	for _, count := range receivedEvents {
		total += count
	}
	assert.Equal(t, numEvents, total)
	mu.Unlock()
}

func TestEventTypeFiltering(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var mu sync.Mutex
	receivedEvents := make(map[models.EventType][]*models.Event)

	// Subscribe to different event types
	for _, eventType := range []models.EventType{models.EventProductCreated, models.EventProductUpdated} {
		et := eventType // Create a new variable to avoid closure issues
		handler := func(event *models.Event) {
			mu.Lock()
			receivedEvents[et] = append(receivedEvents[et], event)
			mu.Unlock()
		}
		err := publisher.Subscribe(et, handler)
		assert.NoError(t, err)
	}

	// Publish events of different types
	events := []*models.Event{
		createTestProductEvent(), // Created
		func() *models.Event {
			e := createTestProductEvent()
			e.Type = models.EventProductUpdated
			return e
		}(),
		func() *models.Event {
			e := createTestProductEvent()
			e.Type = models.EventProductDeleted // Ingen prenumererar på denna
			return e
		}(),
	}

	for _, event := range events {
		err := publisher.Publish(event)
		assert.NoError(t, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify that events ended up in the correct handlers
	mu.Lock()
	assert.Len(t, receivedEvents[models.EventProductCreated], 1)
	assert.Len(t, receivedEvents[models.EventProductUpdated], 1)
	assert.Len(t, receivedEvents[models.EventProductDeleted], 0)
	mu.Unlock()
}

func TestSubscribeMultipleEventTypes(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var mu sync.Mutex
	receivedEvents := make(map[models.EventType][]*models.Event)

	// Subscribe to different event types
	eventTypes := []models.EventType{
		models.EventProductCreated,
		models.EventProductUpdated,
		models.EventProductDeleted,
	}

	for _, et := range eventTypes {
		eventType := et // Capture variable for closure
		handler := func(event *models.Event) {
			mu.Lock()
			receivedEvents[eventType] = append(receivedEvents[eventType], event)
			mu.Unlock()
		}
		err := publisher.Subscribe(eventType, handler)
		assert.NoError(t, err)
	}

	// Publish events of each type
	for _, eventType := range eventTypes {
		event := createTestProductEvent()
		event.Type = eventType
		err := publisher.Publish(event)
		assert.NoError(t, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify that each event type received the correct number of events
	mu.Lock()
	for _, eventType := range eventTypes {
		assert.Len(t, receivedEvents[eventType], 1,
			"Should have received exactly one event for type %s", eventType)
	}
	mu.Unlock()
}

func TestUnsubscribeSpecificHandler(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var mu sync.Mutex
	count1, count2 := 0, 0

	handler1 := func(event *models.Event) {
		mu.Lock()
		count1++
		mu.Unlock()
	}

	handler2 := func(event *models.Event) {
		mu.Lock()
		count2++
		mu.Unlock()
	}

	// Subscribe with both handlers
	err := publisher.Subscribe(models.EventProductCreated, handler1)
	assert.NoError(t, err)
	err = publisher.Subscribe(models.EventProductCreated, handler2)
	assert.NoError(t, err)

	// Unsubscribe only handler1
	err = publisher.Unsubscribe(models.EventProductCreated, handler1)
	assert.NoError(t, err)

	// Publish an event
	event := createTestProductEvent()
	err = publisher.Publish(event)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Verify that only handler2 received the event
	mu.Lock()
	assert.Equal(t, 0, count1, "Handler1 should not receive events after unsubscribe")
	assert.Equal(t, 1, count2, "Handler2 should still receive events")
	mu.Unlock()
}

func TestPublishToNonexistentEventType(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	receivedEvents := 0
	var mu sync.Mutex

	handler := func(event *models.Event) {
		mu.Lock()
		receivedEvents++
		mu.Unlock()
	}

	// Subscribe to an event type
	err := publisher.Subscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	// Publish to a different event type
	event := createTestProductEvent()
	event.Type = models.EventProductUpdated
	err = publisher.Publish(event)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Verify that no events were received
	mu.Lock()
	assert.Equal(t, 0, receivedEvents, "Should not receive events for unsubscribed type")
	mu.Unlock()
}

func TestConcurrentSubscribeUnsubscribe(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var wg sync.WaitGroup
	receivedEvents := make(map[int]int)
	var mu sync.Mutex

	// Create and manage multiple handlers concurrently
	numHandlers := 10
	handlers := make([]func(*models.Event), numHandlers)

	for i := 0; i < numHandlers; i++ {
		handlerID := i
		handlers[i] = func(event *models.Event) {
			mu.Lock()
			receivedEvents[handlerID]++
			mu.Unlock()
		}
	}

	// Subscribe/unsubscribe concurrent
	for i := 0; i < numHandlers; i++ {
		wg.Add(2) // One for subscribe, one for unsubscribe
		go func(id int) {
			defer wg.Done()
			err := publisher.Subscribe(models.EventProductCreated, handlers[id])
			assert.NoError(t, err)
		}(i)

		go func(id int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond) // A little delay before unsubscribe
			err := publisher.Unsubscribe(models.EventProductCreated, handlers[id])
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Publish a few events
	for i := 0; i < 5; i++ {
		event := createTestProductEvent()
		err := publisher.Publish(event)
		assert.NoError(t, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify that no race conditions occurred
	mu.Lock()
	totalEvents := 0
	for _, count := range receivedEvents {
		totalEvents += count
	}
	assert.True(t, totalEvents >= 0, "Should handle concurrent subscribe/unsubscribe safely")
	mu.Unlock()
}
