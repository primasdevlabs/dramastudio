package temporal

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"dramastudio/internal/platform/workflow"
	"dramastudio/internal/platform/workflow/activities"
)

// Builder assembles workers per task queue (§83).
type Builder struct {
	client client.Client
	act    *activities.EpisodeActivities
}

func NewBuilder(c client.Client, act *activities.EpisodeActivities) *Builder {
	return &Builder{client: c, act: act}
}

// NewWorker creates a worker on a queue; callers register additional
// workflows/activities on the returned worker before Run/Start.
func (b *Builder) NewWorker(taskQueue string) worker.Worker {
	return worker.New(b.client, taskQueue, worker.Options{})
}

// BuildAll returns workers for every task queue with the episode
// production workflow and activity set registered.
func (b *Builder) BuildAll() []worker.Worker {
	workers := make([]worker.Worker, 0, len(workflow.AllQueues))
	for _, q := range workflow.AllQueues {
		w := b.NewWorker(q)
		w.RegisterWorkflow(workflow.ProduceEpisode)
		if b.act != nil {
			w.RegisterActivity(b.act)
		}
		workers = append(workers, w)
	}
	return workers
}
