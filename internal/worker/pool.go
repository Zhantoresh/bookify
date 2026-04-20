package worker

import (
	"context"
	"log/slog"
	"sync"
)

type Task func() error

type WorkerPool struct {
	tasks      chan Task
	numWorkers int
	logger     *slog.Logger
}

func NewWorkerPool(numWorkers, queueSize int, logger *slog.Logger) *WorkerPool {
	return &WorkerPool{
		tasks:      make(chan Task, queueSize),
		numWorkers: numWorkers,
		logger:     logger,
	}
}

func (p *WorkerPool) Start(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < p.numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			p.logger.Info("worker_started", "worker_id", workerID)
			for {
				select {
				case <-ctx.Done():
					p.logger.Info("worker_stopping", "worker_id", workerID)
					return
				case task, ok := <-p.tasks:
					if !ok {
						return
					}
					if err := task(); err != nil {
						p.logger.Error("worker_task_failed", "worker_id", workerID, "error", err)
					}
				}
			}
		}(i + 1)
	}
}

func (p *WorkerPool) Submit(task Task) {
	select {
	case p.tasks <- task:
	default:
		p.logger.Warn("worker_queue_full")
	}
}

func (p *WorkerPool) Shutdown() {
	close(p.tasks)
}
