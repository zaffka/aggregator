package domain

import "context"

// Configurer is a wrapping interface to avoid direct access to config struct.
type Configurer interface {
	Generators() []Generator
	Aggregators() []Aggregator
	QueueLen() int
	Storage() int
	QueueTopics() []string
}

// QueueHandler hides realization of the queue.
type QueueHandler interface {
	Pub(ctx context.Context, topic string, value interface{})
	Sub(topic string) <-chan interface{}
}

// StorageHandler - storage interface, not realized yet.
type StorageHandler interface{}

// ExecBuilder - this interface is masking assembling process for aggregators and generators.
// At build stage it creates a slice of Executors - runners for the exec.Pool.
type ExecBuilder interface {
	WithQueue(QueueHandler) ExecBuilder
	WithStorage(StorageHandler) ExecBuilder
	Build() []Executor
}

// Executor interface describes any entity to be executed inside of an exec.Pool.
type Executor interface {
	Name() string
	Do(context.Context) error
}
