package internal

import "github.com/zaffka/aggregator/domain"

type ExecBuilderG struct {
	Generators []domain.Generator
	Queue      domain.QueueHandler
	Storage    domain.StorageHandler
}

func (a *ExecBuilderG) WithQueue(qh domain.QueueHandler) domain.ExecBuilder {
	a.Queue = qh

	return a
}

func (a *ExecBuilderG) WithStorage(sh domain.StorageHandler) domain.ExecBuilder {
	return a
}

func (a *ExecBuilderG) Build() []domain.Executor {
	res := make([]domain.Executor, 0, len(a.Generators))
	for i := range a.Generators {
		res = append(res, &generator{&a.Generators[i], a.Queue})
	}

	return res
}
