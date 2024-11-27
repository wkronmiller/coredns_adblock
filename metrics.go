package coredns_adblock

import (
	"sync"

	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "adblock",
	Name:      "request_count_total",
	Help:      "Counter of requests made.",
}, []string{"server"})

var blockedRequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "adblock",
	Name:      "request_count_blocked",
	Help:      "Counter of requests that were blocked.",
}, []string{"server"})

var once sync.Once
