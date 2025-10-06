package orchestrator

import (
	"context"
	"sync"

	"github.com/LSariol/LightHouse/internal/builder"
	"github.com/LSariol/LightHouse/internal/models"
	"github.com/LSariol/LightHouse/internal/watcher"
)

type Orchestrator struct {
	ctx     context.Context
	cancel  context.CancelFunc
	watcher *watcher.Watcher
	builder *builder.Builder
	jobs    chan Job
	results chan Result
	wg      sync.WaitGroup
}

func NewOrchestrator(cfg ConfigDeps) (*Orchestrator, error) {

	return &Orchestrator{
		ctx:     cfg.Context,
		cancel:  cfg.Cancel,
		watcher: cfg.Watcher,
	}, nil
}

func (o *Orchestrator) Start() {
	o.watcher.Start(o.jobs)
	o.spawnWorkers() // This will create new threads
	go o.handleResults()
}
func (o *Orchestrator) Enqueue(repo models.WatchedRepo) {

}

func (o *Orchestrator) Shutdown() {

}
