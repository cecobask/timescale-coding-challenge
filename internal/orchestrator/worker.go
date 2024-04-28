package orchestrator

import (
	"context"
	"time"

	"github.com/cecobask/timescale-coding-challenge/internal/database"
)

type worker struct {
	id    string
	db    *database.Database
	queue chan job
}

type job struct {
	hostname  string
	startTime string
	endTime   string
}

func newWorker(id string, db *database.Database) *worker {
	return &worker{
		id:    id,
		db:    db,
		queue: make(chan job),
	}
}

func (w *worker) enqueueJob(hostname, startTime, endTime string) {
	w.queue <- job{
		hostname:  hostname,
		startTime: startTime,
		endTime:   endTime,
	}
}

func (w *worker) consumeJobs(ctx context.Context, benchmarks chan<- time.Duration, errors chan<- error) {
	for {
		select {
		case j, ok := <-w.queue:
			if !ok {
				return
			}
			bench, err := w.db.BenchmarkQuery(ctx, j.hostname, j.startTime, j.endTime)
			if err != nil {
				errors <- err
			}
			benchmarks <- bench
		default:
			continue
		}
	}
}
