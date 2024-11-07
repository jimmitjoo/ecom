package chaos

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type chaosTransport struct {
	base  http.RoundTripper
	chaos *NetworkChaos
}

func (t *chaosTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Simulate network latency
	time.Sleep(t.chaos.latency)

	// Simulate packet loss
	if rand.Float64() < t.chaos.lossRate {
		return nil, fmt.Errorf("simulated packet loss")
	}

	// If we don't simulate packet loss, use a mock response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}, nil
}
