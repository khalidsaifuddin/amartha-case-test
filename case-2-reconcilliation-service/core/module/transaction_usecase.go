package module

import (
	"context"

	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/entity"
)

type TransactionUsecase interface {
}

type transactionUsecase struct {
	cfg config.Config
}

func NewTransactionUsecase(cfg config.Config) TransactionUsecase {
	return &transactionUsecase{
		cfg: cfg,
	}
}

func (uc *transactionUsecase) ReconcileTransaction(cfg context.Context, request entity.ReconcileTransactionRequest) (err error) {
	return
}
