package worker

import (
	"priority-task-manager/shared/pkg/types"
	"time"
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
	defer p.releaseWorker()

	p.executor.Exec(task)
}

func (p *Pool) releaseWorker() {
	<-p.workers
}

func (p *Pool) WaitForAllDone() {
	for {
		if len(p.workers) == 0 {
			break
		}

		time.Sleep(2 * time.Second)
	}
}
