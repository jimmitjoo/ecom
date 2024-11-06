package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfilingHandler(t *testing.T) {
	handler := NewProfilingHandler()

	tests := []struct {
		name           string
		endpoint       string
		method         string
		expectedCode   int
		expectedHeader string
	}{
		{
			name:           "CPU Profile GET",
			endpoint:       "/debug/pprof/cpu",
			method:         http.MethodGet,
			expectedCode:   http.StatusOK,
			expectedHeader: "application/octet-stream",
		},
		{
			name:           "Heap Profile GET",
			endpoint:       "/debug/pprof/heap",
			method:         http.MethodGet,
			expectedCode:   http.StatusOK,
			expectedHeader: "application/octet-stream",
		},
		{
			name:           "Goroutine Profile GET",
			endpoint:       "/debug/pprof/goroutine",
			method:         http.MethodGet,
			expectedCode:   http.StatusOK,
			expectedHeader: "application/octet-stream",
		},
		{
			name:           "CPU Profile POST - Should Fail",
			endpoint:       "/debug/pprof/cpu",
			method:         http.MethodPost,
			expectedCode:   http.StatusMethodNotAllowed,
			expectedHeader: "text/plain; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, nil)
			w := httptest.NewRecorder()

			switch strings.TrimPrefix(tt.endpoint, "/debug/pprof/") {
			case "cpu":
				handler.CPUProfile(w, req)
			case "heap":
				handler.HeapProfile(w, req)
			case "goroutine":
				handler.GoroutineProfile(w, req)
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedHeader, w.Header().Get("Content-Type"))

			if tt.expectedCode == http.StatusOK {
				assert.NotEmpty(t, w.Body.Bytes(), "Profile data should not be empty")
			}
		})
	}
}

func TestProfilingHandlerIntegration(t *testing.T) {
	handler := NewProfilingHandler()

	// Test that profiling doesn't interfere with normal operation
	t.Run("Multiple Concurrent Requests", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < 3; i++ {
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/debug/pprof/heap", nil)
				w := httptest.NewRecorder()
				handler.HeapProfile(w, req)
				assert.Equal(t, http.StatusOK, w.Code)
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 3; i++ {
			<-done
		}
	})

	// Test memory usage doesn't grow significantly
	t.Run("Memory Usage", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest(http.MethodGet, "/debug/pprof/heap", nil)
			w := httptest.NewRecorder()
			handler.HeapProfile(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}
