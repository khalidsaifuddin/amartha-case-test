package loanrepository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/entity"
	repository_intf "github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/repository"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/pkg/helper"
	"gorm.io/gorm"
)

type repository struct {
	cfg   config.Config
	db    *gorm.DB
	redis *redis.Pool
}

func New(cfg config.Config, db *gorm.DB, redis *redis.Pool) repository_intf.LoanRepository {
	return &repository{
		cfg:   cfg,
		db:    db,
		redis: redis,
	}
}

func (r *repository) GetLoanBySerial(ctx context.Context, serial string) (resp entity.Loan, err error) {
	redisKey := fmt.Sprintf("Loan:%v", serial)

	return helper.GetOrSetCache(r.redis, redisKey, time.Duration(r.cfg.DefaultTTL)*time.Minute, r.cfg.EnableRedisCache, func() (entity.Loan, error) {
		db := r.db.Model(&Loan{})

		db.Where("serial = ?", serial)

		result := Loan{}
		if err := db.First(&result).Error; err != nil {
			return resp, err
		}

		resp = result.ToEntity()

		return resp, nil
	})
}

func (r *repository) GetLoanBySerials(ctx context.Context, serials []string) (resp []entity.Loan, err error) {
	db := r.db.Model(&Loan{})

	db.Where("serial IN ?", serials)

	results := []Loan{}
	if err := db.Find(&results).Error; err != nil {
		return resp, err
	}

	for _, result := range results {
		resp = append(resp, result.ToEntity())
	}

	return resp, nil
}

func (r *repository) GetLoanRepayments(ctx context.Context, request entity.GetLoanRepaymentsRequest) (resp []entity.LoanRepayment, err error) {
	db := r.db.Model(&LoanRepayment{})

	if request.LoanSerial != "" {
		db.Where("loan_serial = ?", request.LoanSerial)
	}

	if len(request.Serials) > 0 {
		db.Where("serial IN ?", request.Serials)
	}

	if !request.MinDate.IsZero() {
		db.Where("repayment_due_date >= ?", request.MinDate.Format(entity.DateFormat))
	}

	if !request.MaxDate.IsZero() {
		db.Where("repayment_due_date <= ?", request.MaxDate.Format(entity.DateFormat))
	}

	if len(request.Statuses) > 0 {
		db.Where("status IN ?", request.Statuses)
	}

	results := []LoanRepayment{}
	if err := db.Order("repayment_due_date DESC").Find(&results).Error; err != nil {
		return resp, err
	}

	for _, result := range results {
		resp = append(resp, result.ToEntity())
	}

	return resp, nil
}

func (r *repository) UpsertRepaymentTransaction(ctx context.Context, request entity.UpsertRepaymentTransactionRequest) (resp entity.MakeRepaymentResponse, err error) {
	// create transaction for atomicity
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		newLoanRepaymentHistory := LoanRepaymentHistory{}
		newLoanRepaymentHistory.FromEntity(request.LoanRepaymentHistory)

		// get principal only repayment and populate serials in one loop
		principalOnlyAmount := 0.0
		relatedSerials := make([]string, 0, len(request.LoanRepayments))

		for _, loanRepayment := range request.LoanRepayments {
			// Append the serial to the new slice
			relatedSerials = append(relatedSerials, loanRepayment.Serial)

			// Calculate principal-only repayment
			principalOnlyAmount += loanRepayment.RepaymentAmount - loanRepayment.InterestAmount
		}

		// Only assign RelatedLoanRepaymentSerials if it's empty
		if len(newLoanRepaymentHistory.RelatedLoanRepaymentSerials) < 1 {
			newLoanRepaymentHistory.RelatedLoanRepaymentSerials = relatedSerials
		}

		if newLoanRepaymentHistory.Serial != "" {
			// update existing record
			if err := tx.Model(&LoanRepaymentHistory{}).Where("serial = ?", newLoanRepaymentHistory.Serial).Updates(newLoanRepaymentHistory).Error; err != nil {
				return err
			}
		} else {
			// create new record
			if err := tx.Model(&LoanRepaymentHistory{}).Create(&newLoanRepaymentHistory).Error; err != nil {
				return err
			}
		}

		// create transaction for each loan repayment
		for _, loanRepayment := range request.LoanRepayments {
			loadRepaymentRecord := LoanRepayment{}
			loadRepaymentRecord.FromEntity(loanRepayment)

			// inject loan repayment serial to record
			loadRepaymentRecord.LoanRepaymentHistorySerial = newLoanRepaymentHistory.Serial

			// update loan repayment record
			if err := tx.Model(&LoanRepayment{}).Where("serial = ?", loadRepaymentRecord.Serial).Updates(&loadRepaymentRecord).Error; err != nil {
				return err
			}
		}

		// create calculate outstanding and principal only amount only if status is paid
		// while for status pending, we do not change the outstanding yet as it is still waiting for borrower to pay the repayment
		if strings.EqualFold(request.LoanRepaymentHistory.Status, entity.LoanRepaymentStatusPaid) {
			existingLoanRecord := request.Loan
			existingLoanRecord.LoanOutstanding = existingLoanRecord.LoanOutstanding - request.LoanRepaymentHistory.RepaymentAmount
			existingLoanRecord.LoanPrincipalOnly = existingLoanRecord.LoanPrincipalOnly - principalOnlyAmount

			// update loan record
			updatedLoanRecord := Loan{}
			updatedLoanRecord.FromEntity(existingLoanRecord)

			if err := tx.Model(&Loan{}).Where("serial = ?", updatedLoanRecord.Serial).Updates(&updatedLoanRecord).Error; err != nil {
				return err
			}
		}

		// revalidate caches
		redisKey := fmt.Sprintf("Loan:%v", request.Loan.Serial)
		go helper.DeleteRedisCache(r.redis, redisKey)

		return nil
	}); err != nil {
		return resp, err
	}

	return
}

func (r *repository) GetLoanRepaymentHistory(ctx context.Context, request entity.GetLoanRepaymentHistoryRequest) (resp entity.LoanRepaymentHistory, err error) {
	db := r.db.Model(&LoanRepaymentHistory{})

	if request.Serial != "" {
		db.Where("serial = ?", request.Serial)
	}

	if request.LoanSerial != "" {
		db.Where("loan_serial = ?", request.LoanSerial)
	}

	if request.Status != "" {
		db.Where("status = ?", request.Status)
	}

	result := LoanRepaymentHistory{}
	if err := db.First(&result).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, err
		}

		return resp, nil
	}

	return result.ToEntity(), nil
}

func (r *repository) GetRepaymentList(ctx context.Context, request entity.GetRepaymentListRequest) (resp entity.GetRepaymentListResponse, err error) {
	db := r.db.Model(&LoanRepayment{})

	if request.LoanSerial != "" {
		db.Where("loan_serial = ?", request.LoanSerial)
	}

	if !request.StartDate.IsZero() {
		db.Where("repayment_due_date >= ?", request.StartDate.Format(entity.DateFormat))
	}

	if !request.EndDate.IsZero() {
		db.Where("repayment_due_date <= ?", request.EndDate.Format(entity.DateFormat))
	}

	if len(request.Statuses) > 0 {
		db.Where("status IN ?", request.Statuses)
	}

	// count total data
	if err := db.Count(&resp.TotalData).Error; err != nil {
		return resp, err
	}

	limit, start, _, _ := helper.SetupListParameter(resp.Page, resp.PageSize)

	results := []LoanRepayment{}
	if err := db.Limit(int(limit)).Offset(int(start)).Order("repayment_due_date ASC").Find(&results).Error; err != nil {
		return resp, err
	}

	resp.Page = request.Page
	resp.PageSize = request.PageSize
	resp.TotalPage = helper.GenerateTotalPage(resp.TotalData, int64(limit))

	for _, result := range results {
		resp.Items = append(resp.Items, result.ToEntity())
	}

	return resp, nil
}

func (r *repository) GetRepaymentHistoryList(ctx context.Context, request entity.GetRepaymentHistoryListRequest) (resp entity.GetRepaymentHistoryListResponse, err error) {
	db := r.db.Model(&LoanRepaymentHistory{})

	if request.LoanSerial != "" {
		db.Where("loan_serial = ?", request.LoanSerial)
	}

	if !request.StartDate.IsZero() {
		db.Where("repayment_time >= ?", request.StartDate.Format(entity.DateFormat))
	}

	if !request.EndDate.IsZero() {
		db.Where("repayment_time <= ?", request.EndDate.Format(entity.DateFormat))
	}

	if len(request.Statuses) > 0 {
		db.Where("status IN ?", request.Statuses)
	}

	// count total data
	if err := db.Count(&resp.TotalData).Error; err != nil {
		return resp, err
	}

	limit, start, _, _ := helper.SetupListParameter(resp.Page, resp.PageSize)

	results := []LoanRepaymentHistory{}
	if err := db.Limit(int(limit)).Offset(int(start)).Order("repayment_time DESC").Find(&results).Error; err != nil {
		return resp, err
	}

	resp.Page = request.Page
	resp.PageSize = request.PageSize
	resp.TotalPage = helper.GenerateTotalPage(resp.TotalData, int64(limit))

	for _, result := range results {
		resp.Items = append(resp.Items, result.ToEntity())
	}

	return resp, nil
}
