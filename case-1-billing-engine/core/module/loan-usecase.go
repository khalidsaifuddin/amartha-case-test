package module

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/entity"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/repository"
	"golang.org/x/sync/errgroup"
)

type LoanUsecase interface {
	UpsertRepaymentTransaction(ctx context.Context, request entity.UpsertRepaymentTransactionRequest) (resp entity.MakeRepaymentResponse, err error)
	MakeRepayment(ctx context.Context, request entity.MakeRepaymentRequest) (resp entity.MakeRepaymentResponse, err error)
	IsDelinquentByLoanSerial(ctx context.Context, request entity.IsDelinquentByLoanSerialRequest) (resp entity.IsDelinquentByLoanSerialResponse, err error)
	GetOutstandingByLoanSerial(ctx context.Context, request entity.GetOutstandingByLoanSerialRequest) (resp entity.Loan, err error)
	GetRepaymentList(ctx context.Context, request entity.GetRepaymentListRequest) (resp entity.GetRepaymentListResponse, err error)
	GetNextRepayment(ctx context.Context, request entity.GetUnpaidLoanRepaymentRequest) (resp []entity.LoanRepayment, err error)
	GetRepaymentHistoryList(ctx context.Context, request entity.GetRepaymentHistoryListRequest) (resp entity.GetRepaymentHistoryListResponse, err error)
}

type loanUsecase struct {
	cfg      config.Config
	loanRepo repository.LoanRepository
}

func NewLoanUsecase(cfg config.Config, loanRepo repository.LoanRepository) LoanUsecase {
	return &loanUsecase{
		cfg:      cfg,
		loanRepo: loanRepo,
	}
}

func (uc *loanUsecase) UpsertRepaymentTransaction(ctx context.Context, request entity.UpsertRepaymentTransactionRequest) (resp entity.MakeRepaymentResponse, err error) {
	// get loan record
	loan, err := uc.loanRepo.GetLoanBySerial(ctx, request.Loan.Serial)
	if err != nil {
		return resp, err
	}

	if loan.Serial == "" {
		return resp, entity.ErrLoanNotFound
	}

	request.Loan = loan

	// validate loan repayments
	loanRepaymentMap := make(map[string]entity.LoanRepayment)
	loanRepaymentSerials := []string{}

	for i, loanRepayment := range request.LoanRepayments {
		loanRepaymentMap[loanRepayment.Serial] = loanRepayment
		loanRepaymentSerials = append(loanRepaymentSerials, loanRepayment.Serial)

		// impose repayment time to sync with repayment time history
		request.LoanRepayments[i].RepaymentTime = request.LoanRepaymentHistory.RepaymentTime

		// impose repayment status to be equal to history if status is not canceled
		if !strings.EqualFold(request.LoanRepaymentHistory.Status, entity.LoanRepaymentStatusCanceled) {
			// set status to equal with loan repayment history
			request.LoanRepayments[i].Status = request.LoanRepaymentHistory.Status
		} else {
			// set status back to unpaid
			request.LoanRepayments[i].Status = entity.LoanRepaymentStatusUnpaid
		}
	}

	if len(loanRepaymentSerials) < 1 {
		return resp, entity.ErrLoanRepaymentRecordRequired
	}

	existingLoanRepayments, err := uc.loanRepo.GetLoanRepayments(ctx, entity.GetLoanRepaymentsRequest{
		Serials: loanRepaymentSerials,
	})
	if err != nil {
		return resp, fmt.Errorf("error in GetLoanRepaymentBySerials: %v", err)
	}

	// iterate through existingLoanRepayments to check if repayment amount equals to designated amout
	repaymentAmountInvalid := 0
	for _, existingLoanRepayment := range existingLoanRepayments {
		if existingLoanRepayment.TotalAmount != loanRepaymentMap[existingLoanRepayment.Serial].RepaymentAmount {
			repaymentAmountInvalid++
			resp.LoanRepayments = append(resp.LoanRepayments, loanRepaymentMap[existingLoanRepayment.Serial])
		}
	}

	// validate number of invalid repayment amount if greater than 0
	if repaymentAmountInvalid > 0 {
		return resp, entity.ErrLoanReyamentAmountInvalid
	}

	resp, err = uc.loanRepo.UpsertRepaymentTransaction(ctx, request)
	if err != nil {
		return resp, err
	}

	// retrieving loan repayment back for proofing
	loanRepayments, err := uc.loanRepo.GetLoanRepayments(ctx, entity.GetLoanRepaymentsRequest{
		Serials: loanRepaymentSerials,
	})
	if err != nil {
		// non breaking error, just logging
		log.Printf("error in retrieving loan repayments records back")
	}

	resp.RepaymentTransactionStatus = true
	resp.LoanRepayments = loanRepayments

	return resp, nil
}

func (uc *loanUsecase) GetUnpaidLoanRepayment(ctx context.Context, request entity.GetUnpaidLoanRepaymentRequest) (resp []entity.LoanRepayment, err error) {
	// set default date now if no value is available
	if request.Date.IsZero() {
		request.Date = time.Now()
	}

	// get all unpaid and pending repayment from this loan record
	loanRepaymentRequest := entity.GetLoanRepaymentsRequest{
		LoanSerial: request.LoanSerial,
		MaxDate:    request.Date,
		Statuses:   []string{entity.LoanRepaymentStatusPending, entity.LoanRepaymentStatusUnpaid}, // set to true to get pending and unpaid repayments
	}

	loanRepayments, err := uc.loanRepo.GetLoanRepayments(ctx, loanRepaymentRequest)
	if err != nil {
		return resp, fmt.Errorf("error in GetLoanRepaymentBySerial: %v", err)
	}

	return loanRepayments, nil
}

func (uc *loanUsecase) MakeRepayment(ctx context.Context, request entity.MakeRepaymentRequest) (resp entity.MakeRepaymentResponse, err error) {
	// this method calculates automaticaly the amount of repayment a borrower should pay based on loan serial and transaction date

	loanRepayments, err := uc.GetUnpaidLoanRepayment(ctx, entity.GetUnpaidLoanRepaymentRequest{
		Date:       request.Date,
		LoanSerial: request.LoanSerial,
	})
	if err != nil {
		return resp, err
	}

	if len(loanRepayments) < 1 {
		return resp, entity.ErrNoUnpaidRepayment
	}

	repaymentAmount := 0.0

	for i, loanRepayment := range loanRepayments {
		loanRepayments[i].RepaymentAmount = loanRepayment.TotalAmount
		loanRepayments[i].RepaymentTime = request.Date
		loanRepayments[i].Status = request.Status

		repaymentAmount += loanRepayment.TotalAmount
	}

	// get loan record
	loan, err := uc.loanRepo.GetLoanBySerial(ctx, request.LoanSerial)
	if err != nil {
		return resp, err
	}

	// create new default loan repayment history
	loanRepaymentHistory := entity.LoanRepaymentHistory{
		Loan: entity.Loan{
			Serial: request.LoanSerial,
		},
		LoanRepayment: entity.LoanRepayment{
			Serial: loanRepayments[0].Serial, // take first record of loan payments as reference
		},
		RepaymentAmount:       repaymentAmount,
		RepaymentTime:         request.Date,
		RelatedLoanRepayments: loanRepayments,
		Status:                request.Status,
	}

	existingLoanRepaymentHistory, err := uc.loanRepo.GetLoanRepaymentHistory(ctx, entity.GetLoanRepaymentHistoryRequest{
		LoanSerial: request.LoanSerial,
		Status:     entity.LoanRepaymentStatusPending,
	})
	if err != nil {
		return resp, err
	}

	// replace with existing loan repayment history if status is paid so it will update status instead of creating new record
	if strings.EqualFold(request.Status, entity.LoanRepaymentStatusPaid) && existingLoanRepaymentHistory.Serial != "" {
		existingLoanRepaymentHistory.Status = request.Status
		loanRepaymentHistory = existingLoanRepaymentHistory
	}

	// return error if status is not paid but there is existing loan repayment status
	if !strings.EqualFold(request.Status, entity.LoanRepaymentStatusPaid) && existingLoanRepaymentHistory.Serial != "" {
		return resp, fmt.Errorf("there is already pending repayment for this loan")
	}

	transactionRequest := entity.UpsertRepaymentTransactionRequest{
		Loan:                 loan,
		LoanRepayments:       loanRepayments,
		LoanRepaymentHistory: loanRepaymentHistory,
	}

	resp, err = uc.UpsertRepaymentTransaction(ctx, transactionRequest)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (uc *loanUsecase) IsDelinquentByLoanSerial(ctx context.Context, request entity.IsDelinquentByLoanSerialRequest) (resp entity.IsDelinquentByLoanSerialResponse, err error) {
	// get existing loan payment until current date
	if request.MaxDate.IsZero() {
		request.MaxDate = time.Now()
	}

	loanRepaymentRequest := entity.GetLoanRepaymentsRequest{
		LoanSerial: request.LoanSerial,
		MaxDate:    request.MaxDate,
		Statuses:   []string{entity.LoanRepaymentStatusPending, entity.LoanRepaymentStatusUnpaid}, // set to true to get pending and unpaid repayments
	}

	loanRepayments, err := uc.GetLoanRepayments(ctx, loanRepaymentRequest)
	if err != nil {
		return resp, fmt.Errorf("error in GetLoanRepaymentBySerial: %v", err)
	}

	// check if there are any unpaid repayments
	if len(loanRepayments) > uc.cfg.IsDelinquentLoanLimit {
		resp.IsDelinquent = true
	}

	totalUnpaidRepayment := 0.0
	for _, loanRepayment := range loanRepayments {
		totalUnpaidRepayment += loanRepayment.TotalAmount
	}

	resp.UnpaidLoanRepayments = loanRepayments
	resp.TotalUnpaidAmount = totalUnpaidRepayment

	return resp, nil
}

func (uc *loanUsecase) GetOutstandingByLoanSerial(ctx context.Context, request entity.GetOutstandingByLoanSerialRequest) (resp entity.Loan, err error) {
	loan, err := uc.loanRepo.GetLoanBySerial(ctx, request.LoanSerial)
	if err != nil {
		return resp, err
	}

	return loan, nil
}

func (uc *loanUsecase) GetLoanRepayments(ctx context.Context, request entity.GetLoanRepaymentsRequest) (resp []entity.LoanRepayment, err error) {
	resp, err = uc.loanRepo.GetLoanRepayments(ctx, request)
	if err != nil {
		return resp, err
	}

	loanSerials := []string{}
	for _, record := range resp {
		loanSerials = append(loanSerials, record.Loan.Serial)
	}

	loans, err := uc.loanRepo.GetLoanBySerials(ctx, loanSerials)
	if err != nil {
		return resp, err
	}

	loanMap := make(map[string]entity.Loan)
	for _, loan := range loans {
		loanMap[loan.Serial] = loan
	}

	for i, record := range resp {
		resp[i].Loan = loanMap[record.Loan.Serial]
	}

	return resp, nil
}

func (uc *loanUsecase) GetRepaymentList(ctx context.Context, request entity.GetRepaymentListRequest) (resp entity.GetRepaymentListResponse, err error) {
	resp, err = uc.loanRepo.GetRepaymentList(ctx, request)
	if err != nil {
		return resp, err
	}

	loanSerials := []string{}
	for _, item := range resp.Items {
		loanSerials = append(loanSerials, item.Loan.Serial)
	}

	loans, err := uc.loanRepo.GetLoanBySerials(ctx, loanSerials)
	if err != nil {
		return resp, err
	}

	loanMap := make(map[string]entity.Loan)
	for _, loan := range loans {
		loanMap[loan.Serial] = loan
	}

	for i, item := range resp.Items {
		resp.Items[i].Loan = loanMap[item.Loan.Serial]
	}

	return resp, nil
}

func (uc *loanUsecase) GetNextRepayment(ctx context.Context, request entity.GetUnpaidLoanRepaymentRequest) (resp []entity.LoanRepayment, err error) {
	nextRepayment, err := uc.GetUnpaidLoanRepayment(ctx, request)
	if err != nil {
		return resp, err
	}

	return nextRepayment, nil
}

func (uc *loanUsecase) GetRepaymentHistoryList(ctx context.Context, request entity.GetRepaymentHistoryListRequest) (resp entity.GetRepaymentHistoryListResponse, err error) {
	resp, err = uc.loanRepo.GetRepaymentHistoryList(ctx, request)
	if err != nil {
		return resp, err
	}

	loanSerials := []string{}
	loanRepaymentSerials := []string{}

	for _, item := range resp.Items {
		loanSerials = append(loanSerials, item.Loan.Serial)
		loanRepaymentSerials = append(loanRepaymentSerials, item.LoanRepayment.Serial)
	}

	// get data using errgroup
	eg, _ := errgroup.WithContext(ctx)

	loanMapChan := make(chan map[string]entity.Loan)
	eg.Go(func() error {
		loans, err := uc.loanRepo.GetLoanBySerials(ctx, loanSerials)
		if err != nil {
			return err
		}

		loanMap := make(map[string]entity.Loan)
		for _, loan := range loans {
			loanMap[loan.Serial] = loan
		}

		loanMapChan <- loanMap
		return nil
	})

	loanRepaymentMapChan := make(chan map[string]entity.LoanRepayment)
	eg.Go(func() error {
		loanRepayments, err := uc.loanRepo.GetLoanRepayments(ctx, entity.GetLoanRepaymentsRequest{
			Serials: loanRepaymentSerials,
		})
		if err != nil {
			return err
		}

		loanRepaymentMap := make(map[string]entity.LoanRepayment)
		for _, loanRepayment := range loanRepayments {
			loanRepaymentMap[loanRepayment.Serial] = loanRepayment
		}

		loanRepaymentMapChan <- loanRepaymentMap
		return nil
	})

	loanMap := <-loanMapChan
	loanRepaymentMap := <-loanRepaymentMapChan

	if err := eg.Wait(); err != nil {
		return resp, err
	}

	for i, item := range resp.Items {
		resp.Items[i].Loan = loanMap[item.Loan.Serial]
		resp.Items[i].LoanRepayment = loanRepaymentMap[item.LoanRepayment.Serial]
	}

	return resp, nil
}
