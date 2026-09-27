package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dramastudio/internal/platform/workflow"
	"dramastudio/internal/platform/workflow/activities"
)

func main() {
	log.Println("Starting DramaStudio Background Worker...")
	log.Println("Initializing task queues: story-tasks, visual-tasks, video-tasks, audio-tasks, media-tasks, qa-tasks, publishing-tasks...")

	act := activities.NewEpisodeActivities()
	wf := workflow.NewProduceEpisodeWorkflow(act)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Simulate Worker Loop listening for production jobs
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				log.Println("[Worker] Listening on Temporal task queues: video-tasks, story-tasks... (Idle)")
			}
		}
	}()

	_ = wf

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down DramaStudio Worker gracefully...")
}
