package entity

import "time"

const (
	ReconcileTransactionInProgress = "in_progress"
	ReconcileTransactionCompleted  = "completed"
	ReconcileTransactionFailed     = "failed"

	TaskTypeTriggerReconcileTransaction = "trigger_reconcile_transaction"
)

type ReconcileTransactionRequest struct {
	SystemTransactionPath string    `json:"system_transaction_path"`
	BankTransactionPaths  []string  `json:"bank_transaction_paths"`
	StartDateStr          string    `json:"start_date"`
	StartDate             time.Time `json:"-"`
	EndDateStr            string    `json:"end_date"`
	EndDate               time.Time `json:"-"`
	IsAsync               bool      `json:"is_async"`
}

type ReconcileTransactionDetail struct {
	Percentage                 float64                    `json:"percentage"`
	TotalTransaction           int64                      `json:"total_transaction"`
	TotalTranscationProcessed  int64                      `json:"total_transaction_processed"`
	TotalTransactionMatched    int64                      `json:"total_transaction_matched"`
	TotalTransactionUnmatched  int64                      `json:"total_transaction_unmatched"`
	UnmatchedTransactionDetail UnmatchedTransactionDetail `json:"unmatched_transaction_detail"`
}

type UnmatchedTransactionDetail struct {
	UnmatchedSystemTransactions []SystemTransaction          `json:"unmatched_system_transaction"`
	UnmatchedBankTransaction    map[string][]BankTransaction `json:"unmatched_bank_transaction"`
}

type SystemTransaction struct {
	TrxID           string    `json:"trxID"`
	Amount          float64   `json:"amount"`
	Type            string    `json:"type"`
	TransactionTime time.Time `json:"transactionTime"`
}

type BankTransaction struct {
	BankName         string    `json:"bank_name"`
	UniqueIdentifier string    `json:"unique_identifier"`
	Amount           float64   `json:"amount"`
	Date             time.Time `json:"date"`
}

type ReconcileTransactionResponse struct {
	Serial                     string                     `json:"serial"`
	ProgressPercentage         float64                    `json:"progress_percentage"`
	ProgressStatus             string                     `json:"progress_status"`
	ProgressDetailURL          string                     `json:"progress_detail_url"`
	ReconcileTransactionDetail ReconcileTransactionDetail `json:"reconcile_transaction_detail"`
}

type ReconcileTransactionWorkerPayload struct {
	Request  ReconcileTransactionRequest  `json:"request"`
	Progress ReconcileTransactionResponse `json:"progress"`
}
