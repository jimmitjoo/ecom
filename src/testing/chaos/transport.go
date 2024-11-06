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
	// Simulera nätverkslatens
	time.Sleep(t.chaos.latency)

	// Simulera paketförlust
	if rand.Float64() < t.chaos.lossRate {
		return nil, fmt.Errorf("simulerad paketförlust")
	}

	// Om vi inte simulerar paketförlust, använd en mock response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
	}, nil
}
