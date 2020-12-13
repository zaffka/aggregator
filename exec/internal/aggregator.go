package internal

import (
	"context"
	"fmt"

	"github.com/zaffka/aggregator/domain"
)

type aggregator struct {
	*domain.Aggregator
	domain.QueueHandler
	domain.StorageHandler
}

func (a *aggregator) Name() string {
	return a.ID + "." + a.Topic
}

func (a *aggregator) Do(ctx context.Context) error {
	ctx, cancelF := context.WithTimeout(ctx, a.WorkDuration)
	defer cancelF()

	dataCh := a.Sub(a.Topic)

	for {
		select {
		case value := <-dataCh:
			fmt.Printf("> %s=%v\n", a.Topic, value)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
