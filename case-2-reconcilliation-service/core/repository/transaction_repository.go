package repository

import (
	"context"

	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/entity"
)

type TransactionRepository interface {
	SetTransactionReconciliationProgress(updatedProgress entity.ReconcileTransactionResponse) (progress entity.ReconcileTransactionResponse, err error)
	GetTransactionReconciliationProgress(ctx context.Context, serial string) (resp entity.ReconcileTransactionResponse, err error)
}
