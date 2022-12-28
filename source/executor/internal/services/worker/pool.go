package worker

import (
	"priority-task-manager/shared/pkg/types"
)

type Pool struct {
	workers  chan struct{}
	executor Executor
}

func MakePool(maxWorkersCount int, executor Executor) *Pool {
	return &Pool{
		workers:  make(chan struct{}, maxWorkersCount),
		executor: executor,
	}
}

func (p *Pool) Exec(task types.Task) {
	p.workers <- struct{}{}

	go p.exec(task)
}

func (p *Pool) exec(task types.Task) {
	defer func() {
		<-p.workers
	}()

	p.executor.Exec(task)
}
