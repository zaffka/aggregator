package queue

import (
	"context"

	"github.com/zaffka/aggregator/shutdown"
	"go.uber.org/zap"
)

type msgData struct {
	topic string
	value interface{}
}

type listener struct {
	ready chan struct{}
	data  chan interface{}
}

// Queue is the struct holding all the dependencies needed to handle messages.
// You need to use NewQueue() constructor func to init the Queue and run manager routine.
// Manager goroutine is distributing incoming messages for the listening aggregators.
type Queue struct {
	messages       chan msgData
	logger         *zap.SugaredLogger
	topicListeners map[string]listener
}

// NewQueue is a constructor func to init the Queue and run manager routine.
func NewQueue(ctx context.Context, topics []string, queueLength int, log *zap.SugaredLogger) *Queue {
	q := &Queue{
		messages:       make(chan msgData, queueLength),
		logger:         log.Named("queue"),
		topicListeners: make(map[string]listener, len(topics)),
	}

	for _, t := range topics {
		q.topicListeners[t] = listener{
			data:  make(chan interface{}, 1),
			ready: make(chan struct{}),
		}
	}

	go q.manager(ctx)

	return q
}

func (q *Queue) Pub(ctx context.Context, topic string, value interface{}) {
	lisnr := q.topicListeners[topic]
	select {
	case <-lisnr.ready:
	case <-ctx.Done():
		return
	}

	select {
	case q.messages <- msgData{topic, value}:
		q.logger.Named(topic).Debugw("pub func called", "value", value)
	default:
		q.logger.Named(topic).Warnw("dropping the data: queue is full", "value", value)
	}
}

func (q *Queue) Sub(topic string) <-chan interface{} {
	lisnr := q.topicListeners[topic]
	close(lisnr.ready)

	return lisnr.data
}

func (q *Queue) manager(ctx context.Context) {
	shutdown.WaitMe()
	defer shutdown.ImDone()

	l := q.logger.Named("manager")
	l.Info("\u263A starting routine")
	defer l.Info("\u2620 routine stopped")
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-q.messages:
			select {
			case <-ctx.Done():
				return
			case q.topicListeners[msg.topic].data <- msg.value:
			default:
				l.Warn("data dropped: subscriber busy", "topic", msg.topic, "value", msg.value)
			}
		}
	}
}
