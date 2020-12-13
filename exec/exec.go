package exec

import (
	"context"
	"errors"

	"github.com/zaffka/aggregator/domain"
	"github.com/zaffka/aggregator/exec/internal"
	"github.com/zaffka/aggregator/shutdown"
	"go.uber.org/zap"
)

var (
	ErrNoExecBuilder = errors.New("no exec builder for the entity")
)

type Pool struct {
	ctx    context.Context
	logger *zap.SugaredLogger
}

func (p *Pool) Add(execs ...domain.Executor) {
	for _, ex := range execs {
		go func(ex domain.Executor) {
			shutdown.WaitMe()
			defer shutdown.ImDone()

			p.logger.Infow("executing", "name", ex.Name())
			err := ex.Do(p.ctx)
			if err != nil {
				p.logger.Infow("execution finished", "message", err.Error(), "name", ex.Name())
			} else {
				p.logger.Infow("execution finished", "name", ex.Name())
			}
		}(ex)
	}
}

func NewPool(ctx context.Context, logger *zap.SugaredLogger) *Pool {
	p := &Pool{
		ctx:    ctx,
		logger: logger.Named("exec.pool"),
	}

	return p
}

func For(entity interface{}) domain.ExecBuilder {
	switch e := entity.(type) {
	case []domain.Aggregator:
		return &internal.ExecBuilderA{Aggregators: e}
	case []domain.Generator:
		return &internal.ExecBuilderG{Generators: e}
	default:
		return nil
	}
}
