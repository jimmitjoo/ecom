package handlers

import (
	"net/http"
	"runtime"
	"runtime/pprof"
	"time"
)

// ProfilingHandler handles performance profiling endpoints
type ProfilingHandler struct{}

// NewProfilingHandler creates a new profiling handler
func NewProfilingHandler() *ProfilingHandler {
	return &ProfilingHandler{}
}

// CPUProfile handles CPU profiling requests
func (h *ProfilingHandler) CPUProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	pprof.StartCPUProfile(w)
	defer pprof.StopCPUProfile()

	// Run for 30 seconds to collect profile data
	time.Sleep(30 * time.Second)
	runtime.GC()
}

// HeapProfile handles heap memory profiling requests
func (h *ProfilingHandler) HeapProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	runtime.GC() // Run garbage collection before heap dump
	pprof.WriteHeapProfile(w)
}

// GoroutineProfile handles goroutine profiling requests
func (h *ProfilingHandler) GoroutineProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	p := pprof.Lookup("goroutine")
	p.WriteTo(w, 1)
}
