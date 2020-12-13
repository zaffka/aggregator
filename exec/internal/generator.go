package internal

import (
	"context"
	"time"

	"github.com/zaffka/aggregator/domain"
)

type generator struct {
	*domain.Generator
	domain.QueueHandler
}

func (g *generator) Name() string {
	return "generator." + g.Topic
}

func (g *generator) Do(ctx context.Context) error {
	ctx, cancelF := context.WithTimeout(ctx, g.WorkDuration)
	defer cancelF()

	ticker := time.NewTicker(g.GenPeriod)
	defer ticker.Stop()

	v := g.StartValue

	for {
		g.Pub(ctx, g.Topic, v)
		v += getShiftValue(g.MaxStep)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
