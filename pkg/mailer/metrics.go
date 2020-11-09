package mailer

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	parsedPerformers *prometheus.CounterVec
}
