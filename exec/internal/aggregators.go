package internal

import "github.com/zaffka/aggregator/domain"

type ExecBuilderA struct {
	Aggregators []domain.Aggregator
	Queue       domain.QueueHandler
	Storage     domain.StorageHandler
}

func (a *ExecBuilderA) WithQueue(qh domain.QueueHandler) domain.ExecBuilder {
	a.Queue = qh

	return a
}

func (a *ExecBuilderA) WithStorage(sh domain.StorageHandler) domain.ExecBuilder {
	a.Storage = sh

	return a
}

func (a *ExecBuilderA) Build() []domain.Executor {
	res := make([]domain.Executor, 0, len(a.Aggregators))
	for i := range a.Aggregators {
		res = append(res, &aggregator{&a.Aggregators[i], a.Queue, a.Storage})
	}

	return res
}
