package orchestrator

import (
	"context"
	"sync"

	"github.com/LSariol/LightHouse/internal/builder"
)

type ConfigDeps struct {
	Context context.Context
	Cancel  context.CancelFunc
	Watcher *watcher.Watcher
	Builder *builder.Builder
	Jobs    chan Job
	Results chan Result
	Wg      sync.WaitGroup
}
