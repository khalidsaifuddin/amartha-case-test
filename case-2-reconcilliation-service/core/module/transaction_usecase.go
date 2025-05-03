package module

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/entity"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/repository"
	workerclient "github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/handler/worker/client"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/pkg/helper"
)

type TransactionUsecase interface {
	RunReconcileTransaction(ctx context.Context, request entity.ReconcileTransactionRequest) (resp entity.ReconcileTransactionResponse, err error)
	ReconcileTransaction(ctx context.Context, request entity.ReconcileTransactionRequest, progressResponse entity.ReconcileTransactionResponse) (resp entity.ReconcileTransactionResponse, err error)
	GetTransactionReconciliationProgress(ctx context.Context, serial string) (resp entity.ReconcileTransactionResponse, err error)
	SetTransactionReconciliationProgress(ctx context.Context, updatedProgress entity.ReconcileTransactionResponse) (progress entity.ReconcileTransactionResponse, err error)
}

type transactionUsecase struct {
	cfg             config.Config
	transactionRepo repository.TransactionRepository
	workerClient    workerclient.WorkerClient
}

func NewTransactionUsecase(cfg config.Config, transactionRepo repository.TransactionRepository, workerClient workerclient.WorkerClient) TransactionUsecase {
	return &transactionUsecase{
		cfg:             cfg,
		transactionRepo: transactionRepo,
		workerClient:    workerClient,
	}
}

func (uc *transactionUsecase) RunReconcileTransaction(ctx context.Context, request entity.ReconcileTransactionRequest) (resp entity.ReconcileTransactionResponse, err error) {
	// create new transaction reconciliation progress in repository, set default value is in progress
	progress, err := uc.transactionRepo.SetTransactionReconciliationProgress(entity.ReconcileTransactionResponse{})
	if err != nil {
		return resp, err
	}

	if request.IsAsync {
		log.Printf("trigger async task to worker server")
		// trigger async task to worker server
		// create task payload
		payload := entity.ReconcileTransactionWorkerPayload{}
		payload.Request = request
		payload.Progress = progress

		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return resp, err
		}

		_, err = uc.workerClient.Enqueue(entity.TaskTypeTriggerReconcileTransaction, payloadJSON)
		if err != nil {
			return resp, err
		}

		return progress, nil
	}

	progress, err = uc.ReconcileTransaction(ctx, request, progress)
	if err != nil {
		progress.ProgressStatus = entity.ReconcileTransactionFailed

		return progress, err
	}

	return progress, nil
}

func (uc *transactionUsecase) ReconcileTransaction(ctx context.Context, request entity.ReconcileTransactionRequest, progressResponse entity.ReconcileTransactionResponse) (resp entity.ReconcileTransactionResponse, err error) {
	// logging in this function is verbose to track the progress of the transaction reconciliation
	// especially when the transaction is asynchronous

	systemTransactions, err := uc.ConvertSystemTransactionToStruct(ctx, request.SystemTransactionPath)
	if err != nil {
		return resp, err
	}

	bankTransactionMap := make(map[string][]entity.BankTransaction)
	for _, csvPath := range request.BankTransactionPaths {
		bankTransactions, err := uc.ConvertBankTransactionToStruct(ctx, csvPath)
		if err != nil {
			return resp, err
		}

		if len(bankTransactions) < 1 {
			continue
		}

		bankTransactionMap[bankTransactions[0].BankName] = bankTransactions
	}

	// system transaction progress
	log.Printf("system transaction progress, populating hashmap")

	// create hashmap so the process will be O(1) and prevent O(n^2) complexity
	systemTransactionHashMap := make(map[string]entity.SystemTransaction)
	for _, transaction := range systemTransactions {
		amount := transaction.Amount
		if strings.EqualFold(transaction.Type, "debit") {
			amount = -transaction.Amount
		}

		key := fmt.Sprintf("%f-%s", amount, transaction.TransactionTime.Format(entity.DateFormat))
		systemTransactionHashMap[key] = transaction
	}

	// bank transaction progress
	log.Printf("bank transaction progress, populating hashmap")

	bankTransactionHashMap := make(map[string]entity.BankTransaction)
	for _, bankTransactions := range bankTransactionMap {
		for _, transaction := range bankTransactions {
			key := fmt.Sprintf("%f-%s", transaction.Amount, transaction.Date.Format(entity.DateFormat))
			bankTransactionHashMap[key] = transaction
		}
	}

	log.Printf("system transaction progress, processing transactions")

	transactionProgress := entity.ReconcileTransactionDetail{
		Percentage:                0.0,
		TotalTransaction:          int64(len(systemTransactionHashMap)) + int64(len(bankTransactionHashMap)),
		TotalTranscationProcessed: 0,
		TotalTransactionMatched:   0,
		TotalTransactionUnmatched: 0,
		UnmatchedTransactionDetail: entity.UnmatchedTransactionDetail{
			UnmatchedSystemTransactions: make([]entity.SystemTransaction, 0),
			UnmatchedBankTransaction:    make(map[string][]entity.BankTransaction),
		},
	}

	for key, transaction := range systemTransactionHashMap {
		transactionProgress.TotalTranscationProcessed++

		if _, ok := bankTransactionHashMap[key]; ok {
			// matched
			transactionProgress.TotalTransactionMatched++

			log.Printf("MATCHED system transaction key: %s, value: %s", key, transaction.TrxID)
		} else {
			// unmatched
			transactionProgress.TotalTransactionUnmatched++

			transactionProgress.UnmatchedTransactionDetail.UnmatchedSystemTransactions = append(transactionProgress.UnmatchedTransactionDetail.UnmatchedSystemTransactions, transaction)

			log.Printf("UNMATCHED system transaction key: %s, value: %s", key, transaction.TrxID)
		}

		systemTransactionProgressPercentage := float64(transactionProgress.TotalTranscationProcessed) / float64(transactionProgress.TotalTransaction) * 100
		transactionProgress.Percentage = systemTransactionProgressPercentage
	}

	log.Printf("bank transaction progress, processing transactions")

	for key, transaction := range bankTransactionHashMap {
		transactionProgress.TotalTranscationProcessed++

		if _, ok := systemTransactionHashMap[key]; ok {
			// matched
			transactionProgress.TotalTransactionMatched++

			log.Printf("MATCHED bank transaction key: %s, value: %s", key, transaction.UniqueIdentifier)
		} else {
			// unmatched
			transactionProgress.TotalTransactionUnmatched++

			if _, ok := transactionProgress.UnmatchedTransactionDetail.UnmatchedBankTransaction[transaction.BankName]; !ok {
				transactionProgress.UnmatchedTransactionDetail.UnmatchedBankTransaction[transaction.BankName] = make([]entity.BankTransaction, 0)
			}

			transactionProgress.UnmatchedTransactionDetail.UnmatchedBankTransaction[transaction.BankName] = append(transactionProgress.UnmatchedTransactionDetail.UnmatchedBankTransaction[transaction.BankName], transaction)

			log.Printf("UNMATCHED bank transaction key: %s, value: %s", key, transaction.UniqueIdentifier)
		}

		bankTransactionProgressPercentage := float64(transactionProgress.TotalTranscationProcessed) / float64(transactionProgress.TotalTransaction) * 100
		transactionProgress.Percentage = bankTransactionProgressPercentage
	}

	if transactionProgress.Percentage >= 100 {
		progressResponse.ProgressStatus = entity.ReconcileTransactionCompleted
		progressResponse.ProgressPercentage = transactionProgress.Percentage
	}

	progressResponse.ReconcileTransactionDetail = transactionProgress

	return progressResponse, nil
}

func (uc *transactionUsecase) GetTransactionReconciliationProgress(ctx context.Context, serial string) (resp entity.ReconcileTransactionResponse, err error) {
	return uc.transactionRepo.GetTransactionReconciliationProgress(ctx, serial)
}

func (uc *transactionUsecase) SetTransactionReconciliationProgress(ctx context.Context, updatedProgress entity.ReconcileTransactionResponse) (progress entity.ReconcileTransactionResponse, err error) {
	return uc.transactionRepo.SetTransactionReconciliationProgress(updatedProgress)
}

func (uc *transactionUsecase) ConvertSystemTransactionToStruct(ctx context.Context, csvPath string) (resp []entity.SystemTransaction, err error) {
	dataset, err := helper.ParseCSV(csvPath)
	if err != nil {
		return resp, err
	}

	for _, data := range dataset {
		entityRecord := entity.SystemTransaction{
			TrxID: data["trxID"],
			Amount: func() float64 {
				amount, _ := strconv.ParseFloat(data["amount"], 64)
				return amount
			}(),
			Type: data["type"],
			TransactionTime: func() time.Time {
				t, _ := time.Parse("2006-01-02 15:04:05", data["transactionTime"])
				return t
			}(),
		}

		resp = append(resp, entityRecord)
	}

	return resp, nil
}

func (uc *transactionUsecase) ConvertBankTransactionToStruct(ctx context.Context, csvPath string) (resp []entity.BankTransaction, err error) {
	// get bank name from csv path, csv file name, for example: /tmp/bca.csv
	// will be converted to bca
	bankName := csvPath
	bankName = bankName[strings.LastIndex(bankName, "/")+1:]
	bankName = bankName[:strings.LastIndex(bankName, ".")]
	bankName = strings.ToLower(bankName)

	dataset, err := helper.ParseCSV(csvPath)
	if err != nil {
		return resp, err
	}

	for _, data := range dataset {
		entityRecord := entity.BankTransaction{
			BankName:         bankName,
			UniqueIdentifier: data["unique_identifier"],
			Amount: func() float64 {
				amount, _ := strconv.ParseFloat(data["amount"], 64)
				return amount
			}(),
			Date: func() time.Time {
				t, _ := time.Parse("2006-01-02", data["date"])
				return t
			}(),
		}

		resp = append(resp, entityRecord)
	}

	return resp, nil
}
