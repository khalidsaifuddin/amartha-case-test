package worker

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

func (w *WorkerServer) RunDummyTask(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("run dummy task")

	return nil
}
