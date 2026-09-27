package temporal

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Builder assembles workers per task queue.
type Builder struct {
	client client.Client
}

func NewBuilder(c client.Client) *Builder {
	return &Builder{client: c}
}

// NewWorker creates a worker on a queue; callers register workflows and
// activities on the returned worker before Run/Start.
func (b *Builder) NewWorker(taskQueue string) worker.Worker {
	return worker.New(b.client, taskQueue, worker.Options{})
}
