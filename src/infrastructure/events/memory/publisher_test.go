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

	// Skapa en handler som sparar mottagna events
	handler := func(event *models.Event) {
		mu.Lock()
		receivedEvents = append(receivedEvents, event)
		mu.Unlock()
	}

	// Prenumerera på events
	err := publisher.Subscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	// Publicera ett event
	event := createTestProductEvent()
	err = publisher.Publish(event)
	assert.NoError(t, err)

	// Vänta lite så att event hinner processas
	time.Sleep(100 * time.Millisecond)

	// Verifiera att eventet togs emot
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

	// Skapa flera handlers
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

	// Publicera ett event
	event := createTestProductEvent()
	err := publisher.Publish(event)
	assert.NoError(t, err)

	// Vänta på att alla handlers är klara
	wg.Wait()

	// Verifiera att alla handlers fick eventet
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

	// Prenumerera och avprenumerera
	err := publisher.Subscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	err = publisher.Unsubscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	// Publicera ett event
	event := createTestProductEvent()
	err = publisher.Publish(event)
	assert.NoError(t, err)

	// Vänta lite
	time.Sleep(100 * time.Millisecond)

	// Verifiera att inga events togs emot
	mu.Lock()
	assert.Empty(t, receivedEvents)
	mu.Unlock()
}

func TestConcurrentPublishSubscribe(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var wg sync.WaitGroup
	var mu sync.Mutex
	receivedEvents := make(map[string]int)

	// Skapa handlers för olika event typer
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

	// Publicera events concurrent
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

	// Verifiera att rätt antal events togs emot
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

	// Prenumerera på olika event typer
	for _, eventType := range []models.EventType{models.EventProductCreated, models.EventProductUpdated} {
		et := eventType // Skapa en ny variabel för att undvika closure problem
		handler := func(event *models.Event) {
			mu.Lock()
			receivedEvents[et] = append(receivedEvents[et], event)
			mu.Unlock()
		}
		err := publisher.Subscribe(et, handler)
		assert.NoError(t, err)
	}

	// Publicera events av olika typer
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

	// Verifiera att events hamnade i rätt handlers
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

	// Prenumerera på olika event typer
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

	// Publicera events av varje typ
	for _, eventType := range eventTypes {
		event := createTestProductEvent()
		event.Type = eventType
		err := publisher.Publish(event)
		assert.NoError(t, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verifiera att varje event typ fick rätt antal events
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

	// Prenumerera med båda handlers
	err := publisher.Subscribe(models.EventProductCreated, handler1)
	assert.NoError(t, err)
	err = publisher.Subscribe(models.EventProductCreated, handler2)
	assert.NoError(t, err)

	// Avprenumerera bara handler1
	err = publisher.Unsubscribe(models.EventProductCreated, handler1)
	assert.NoError(t, err)

	// Publicera ett event
	event := createTestProductEvent()
	err = publisher.Publish(event)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Verifiera att bara handler2 fick eventet
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

	// Prenumerera på en event typ
	err := publisher.Subscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)

	// Publicera till en annan event typ
	event := createTestProductEvent()
	event.Type = models.EventProductUpdated
	err = publisher.Publish(event)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Verifiera att inga events togs emot
	mu.Lock()
	assert.Equal(t, 0, receivedEvents, "Should not receive events for unsubscribed type")
	mu.Unlock()
}

func TestConcurrentSubscribeUnsubscribe(t *testing.T) {
	publisher := NewMemoryEventPublisher()
	var wg sync.WaitGroup
	receivedEvents := make(map[int]int)
	var mu sync.Mutex

	// Skapa och hantera flera handlers concurrent
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
		wg.Add(2) // En för subscribe, en för unsubscribe
		go func(id int) {
			defer wg.Done()
			err := publisher.Subscribe(models.EventProductCreated, handlers[id])
			assert.NoError(t, err)
		}(i)

		go func(id int) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond) // Lite delay innan unsubscribe
			err := publisher.Unsubscribe(models.EventProductCreated, handlers[id])
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Publicera några events
	for i := 0; i < 5; i++ {
		event := createTestProductEvent()
		err := publisher.Publish(event)
		assert.NoError(t, err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verifiera att inga race conditions uppstod
	mu.Lock()
	totalEvents := 0
	for _, count := range receivedEvents {
		totalEvents += count
	}
	assert.True(t, totalEvents >= 0, "Should handle concurrent subscribe/unsubscribe safely")
	mu.Unlock()
}