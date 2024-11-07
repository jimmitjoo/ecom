package events_test

import (
	"testing"

	"github.com/jimmitjoo/ecom/src/domain/events"
	"github.com/jimmitjoo/ecom/src/domain/models"
	"github.com/stretchr/testify/assert"
)

// MockEventPublisher implements the EventPublisher interface for testing
type MockEventPublisher struct {
	publishCalled   bool
	subscribeCalled bool
	lastEvent       *models.Event
	lastEventType   models.EventType
	lastHandler     func(*models.Event)
}

func (m *MockEventPublisher) Publish(event *models.Event) error {
	m.publishCalled = true
	m.lastEvent = event
	return nil
}

func (m *MockEventPublisher) Subscribe(eventType models.EventType, handler func(*models.Event)) error {
	m.subscribeCalled = true
	m.lastEventType = eventType
	m.lastHandler = handler
	return nil
}

func (m *MockEventPublisher) Unsubscribe(eventType models.EventType, handler func(*models.Event)) error {
	return nil
}

// TestEventPublisherInterface verifies that the interface is implemented correctly
func TestEventPublisherInterface(t *testing.T) {
	var _ events.EventPublisher = &MockEventPublisher{} // Compile-time test
}

// TestEventPublisherUsage verifies that the interface is used correctly
func TestEventPublisherUsage(t *testing.T) {
	publisher := &MockEventPublisher{}

	// Testa Publish
	event := &models.Event{
		ID:   "test_event_1",
		Type: models.EventProductCreated,
		Data: &models.ProductEvent{
			ProductID: "test_product_1",
			Action:    "created",
		},
	}

	err := publisher.Publish(event)
	assert.NoError(t, err)
	assert.True(t, publisher.publishCalled)
	assert.Equal(t, event, publisher.lastEvent)

	// Testa Subscribe
	handler := func(e *models.Event) {}
	err = publisher.Subscribe(models.EventProductCreated, handler)
	assert.NoError(t, err)
	assert.True(t, publisher.subscribeCalled)
	assert.Equal(t, models.EventProductCreated, publisher.lastEventType)
	assert.NotNil(t, publisher.lastHandler)
}

// TestEventTypes verifies that all event types are defined correctly
func TestEventTypes(t *testing.T) {
	eventTypes := []models.EventType{
		models.EventProductCreated,
		models.EventProductUpdated,
		models.EventProductDeleted,
	}

	// Verify that event types are unique
	seen := make(map[models.EventType]bool)
	for _, et := range eventTypes {
		assert.False(t, seen[et], "Event type %s is duplicated", et)
		seen[et] = true
		assert.NotEmpty(t, string(et), "Event type should not be empty")
	}
}

// TestEventData verifies that event data is handled correctly
func TestEventData(t *testing.T) {
	event := &models.Event{
		ID:   "test_event_1",
		Type: models.EventProductCreated,
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

	// Verify that event data can be type asserted correctly
	productEvent, ok := event.Data.(*models.ProductEvent)
	assert.True(t, ok, "Should be able to type assert event data")
	assert.Equal(t, "test_product_1", productEvent.ProductID)
	assert.Equal(t, "created", productEvent.Action)
	assert.NotNil(t, productEvent.Product)
}

// TestEventValidation verifies that event data is validated correctly
func TestEventValidation(t *testing.T) {
	testCases := []struct {
		name      string
		event     *models.Event
		shouldErr bool
	}{
		{
			name: "Valid event",
			event: &models.Event{
				ID:   "test_event_1",
				Type: models.EventProductCreated,
				Data: &models.ProductEvent{
					ProductID: "test_product_1",
					Action:    "created",
					Product: &models.Product{
						ID:        "test_product_1",
						SKU:       "TEST-123",
						BaseTitle: "Test Product",
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "Valid delete event without product",
			event: &models.Event{
				ID:   "test_event_2",
				Type: models.EventProductDeleted,
				Data: &models.ProductEvent{
					ProductID: "test_product_1",
					Action:    "deleted",
				},
			},
			shouldErr: false,
		},
		{
			name: "Missing ID",
			event: &models.Event{
				Type: models.EventProductCreated,
				Data: &models.ProductEvent{
					ProductID: "test_product_1",
					Action:    "created",
					Product: &models.Product{
						ID:        "test_product_1",
						SKU:       "TEST-123",
						BaseTitle: "Test Product",
					},
				},
			},
			shouldErr: true,
		},
		{
			name: "Missing Type",
			event: &models.Event{
				ID: "test_event_1",
				Data: &models.ProductEvent{
					ProductID: "test_product_1",
					Action:    "created",
					Product: &models.Product{
						ID:        "test_product_1",
						SKU:       "TEST-123",
						BaseTitle: "Test Product",
					},
				},
			},
			shouldErr: true,
		},
		{
			name: "Missing Data",
			event: &models.Event{
				ID:   "test_event_1",
				Type: models.EventProductCreated,
			},
			shouldErr: true,
		},
		{
			name: "Missing Product in non-delete event",
			event: &models.Event{
				ID:   "test_event_1",
				Type: models.EventProductCreated,
				Data: &models.ProductEvent{
					ProductID: "test_product_1",
					Action:    "created",
				},
			},
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := models.ValidateEvent(tc.event)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
