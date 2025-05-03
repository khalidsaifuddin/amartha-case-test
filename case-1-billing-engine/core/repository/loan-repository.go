package repository

import (
	"context"

	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/entity"
)

type LoanRepository interface {
	GetLoanRepayments(ctx context.Context, request entity.GetLoanRepaymentsRequest) (resp []entity.LoanRepayment, err error)
	UpsertRepaymentTransaction(ctx context.Context, request entity.UpsertRepaymentTransactionRequest) (resp entity.MakeRepaymentResponse, err error)
	GetLoanBySerial(ctx context.Context, serial string) (resp entity.Loan, err error)
	GetLoanBySerials(ctx context.Context, serials []string) (resp []entity.Loan, err error)
	GetLoanRepaymentHistory(ctx context.Context, request entity.GetLoanRepaymentHistoryRequest) (resp entity.LoanRepaymentHistory, err error)
	GetRepaymentList(ctx context.Context, request entity.GetRepaymentListRequest) (resp entity.GetRepaymentListResponse, err error)
	GetRepaymentHistoryList(ctx context.Context, request entity.GetRepaymentHistoryListRequest) (resp entity.GetRepaymentHistoryListResponse, err error)
}
