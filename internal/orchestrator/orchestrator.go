package orchestrator

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lafikl/consistent"

	"github.com/cecobask/timescale-coding-challenge/internal/database"
	"github.com/cecobask/timescale-coding-challenge/pkg/log"
)

type Orchestrator struct {
	ring    *consistent.Consistent
	workers map[string]*worker
	config  *csv.Reader
}

func New(workerCount int, db *database.Database) *Orchestrator {
	o := &Orchestrator{
		ring:    consistent.New(),
		workers: make(map[string]*worker, workerCount),
	}
	for i := 1; i <= workerCount; i++ {
		id := uuid.NewString()
		o.workers[id] = newWorker(id, db)
		o.ring.Add(id)
	}
	return o
}

func (o *Orchestrator) LoadConfig(path string) error {
	if path == "" {
		o.config = csv.NewReader(os.Stdin)
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed opening csv config file: %w", err)
	}
	o.config = csv.NewReader(file)
	return nil
}

func (o *Orchestrator) Orchestrate(ctx context.Context) ([]time.Duration, error) {
	line, err := o.config.Read()
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("failed to read csv header: %w", err)
		}
	}
	if !slices.Equal(line, csvHeader()) {
		return nil, fmt.Errorf("invalid csv format")
	}
	var (
		waitGroup  = new(sync.WaitGroup)
		benchChan  = make(chan time.Duration)
		errChan    = make(chan error, 1)
		benchmarks = make([]time.Duration, 0)
	)
	for _, currentWorker := range o.workers {
		waitGroup.Add(1)
		go func(w *worker) {
			defer waitGroup.Done()
			w.consumeJobs(ctx, benchChan, errChan)
		}(currentWorker)
	}
	go awaitChannelsOutput(ctx, benchChan, errChan, &benchmarks)
	if err = o.processConfig(); err != nil {
		return nil, err
	}
	waitGroup.Wait()
	return benchmarks, nil
}

func (o *Orchestrator) processConfig() error {
	waitGroup := new(sync.WaitGroup)
	for {
		line, err := o.config.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read csv: %w", err)
		}
		var w *worker
		w, err = o.findWorker(line[0])
		if err != nil {
			return fmt.Errorf("failed to find worker: %w", err)
		}
		waitGroup.Add(1)
		go func(w *worker, line []string) {
			defer waitGroup.Done()
			w.enqueueJob(line[0], line[1], line[2])
		}(w, line)
	}
	waitGroup.Wait()
	for _, w := range o.workers {
		close(w.queue)
	}
	return nil
}

func (o *Orchestrator) findWorker(host string) (*worker, error) {
	workerID, err := o.ring.Get(host)
	if err != nil {
		return nil, fmt.Errorf("could not find any workers in the ring: %w", err)
	}
	w, ok := o.workers[workerID]
	if !ok {
		return nil, fmt.Errorf("could not find worker with id %s", workerID)
	}
	return w, nil
}

func awaitChannelsOutput(ctx context.Context, benchChan <-chan time.Duration, errChan <-chan error, benchmarks *[]time.Duration) {
	for {
		select {
		case bm, ok := <-benchChan:
			if ok {
				*benchmarks = append(*benchmarks, bm)
			}
		case err := <-errChan:
			if err != nil {
				log.FromContext(ctx).ExitOnError(err)
			}
		default:
			continue
		}
	}
}

func csvHeader() []string {
	return []string{
		"hostname",
		"start_time",
		"end_time",
	}
}
