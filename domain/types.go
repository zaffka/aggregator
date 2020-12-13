package domain

import (
	"time"
)

// Aggregator consumes data form the queue.
type Aggregator struct {
	ID           string
	Topic        string
	WorkDuration time.Duration
}

// Generator puts data to the queue.
type Generator struct {
	Topic        string
	StartValue   int
	MaxStep      int
	GenPeriod    time.Duration
	WorkDuration time.Duration
}
