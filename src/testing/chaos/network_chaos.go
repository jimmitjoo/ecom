package chaos

import (
	"net/http"
	"sync"
	"time"
)

// NetworkChaos simulerar nätverksproblem
type NetworkChaos struct {
	latency  time.Duration
	lossRate float64
	mu       sync.Mutex
}

func NewNetworkChaos(latency time.Duration, lossRate float64) *NetworkChaos {
	return &NetworkChaos{
		latency:  latency,
		lossRate: lossRate,
	}
}

func (nc *NetworkChaos) WrapClient(client *http.Client) *http.Client {
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	client.Transport = &chaosTransport{
		base:  transport,
		chaos: nc,
	}
	return client
}
