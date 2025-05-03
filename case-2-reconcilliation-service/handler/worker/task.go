package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/entity"
)

func (w *WorkerServer) RunDummyTask(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("run dummy task")

	return nil
}

func (w *WorkerServer) TriggerReconcileTransaction(ctx context.Context, t *asynq.Task) error {
	fmt.Printf("trigger reconcile transaction")

	payload := entity.ReconcileTransactionWorkerPayload{}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// call usecase
	resp, err := w.transactionUc.ReconcileTransaction(ctx, payload.Request, payload.Progress)
	if err != nil {
		return fmt.Errorf("ReconcileTransaction failed: %v: %w", err, asynq.SkipRetry)
	}

	// set progress to redis
	if _, err := w.transactionUc.SetTransactionReconciliationProgress(ctx, resp); err != nil {
		return fmt.Errorf("SetTransactionReconciliationProgress failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}
