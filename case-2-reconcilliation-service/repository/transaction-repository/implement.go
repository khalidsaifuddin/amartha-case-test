package transactionrepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/entity"
	repository_intf "github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/repository"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/pkg/helper"
)

type transactionRepository struct {
	cfg                                  config.Config
	redis                                *redis.Pool
	transactionReconciliationProgressKey string
}

func New(cfg config.Config, redis *redis.Pool) repository_intf.TransactionRepository {
	return &transactionRepository{
		cfg:                                  cfg,
		redis:                                redis,
		transactionReconciliationProgressKey: "TransactionReconciliationProgress",
	}
}

func (r *transactionRepository) SetTransactionReconciliationProgress(updatedProgress entity.ReconcileTransactionResponse) (progress entity.ReconcileTransactionResponse, err error) {
	if updatedProgress.Serial == "" {
		serial := helper.GenerateUUID()

		progress = entity.ReconcileTransactionResponse{
			Serial:                     serial,
			ProgressPercentage:         0,
			ProgressStatus:             entity.ReconcileTransactionInProgress,
			ProgressDetailURL:          fmt.Sprintf(r.cfg.TransactionReconciliationProgressURL, serial),
			ReconcileTransactionDetail: entity.ReconcileTransactionDetail{},
		}

		updatedProgress = progress
	}

	// convert progress to json
	progressJSON, err := json.Marshal(updatedProgress)
	if err != nil {
		return progress, err
	}

	// set progress to redis
	redisKey := fmt.Sprintf("%v:%v", r.transactionReconciliationProgressKey, updatedProgress.Serial)
	if err := helper.SetRedisCache(r.redis, redisKey, progressJSON, 0); err != nil {
		return updatedProgress, err
	}

	return updatedProgress, nil
}

func (r *transactionRepository) GetTransactionReconciliationProgress(ctx context.Context, serial string) (resp entity.ReconcileTransactionResponse, err error) {
	redisKey := fmt.Sprintf("%v:%v", r.transactionReconciliationProgressKey, serial)
	cache, err := helper.GetRedisCache(r.redis, redisKey)
	if err != nil {
		return resp, err
	}

	if cache == "" {
		return resp, fmt.Errorf("transaction reconciliation progress not found")
	}

	err = json.Unmarshal([]byte(cache), &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
