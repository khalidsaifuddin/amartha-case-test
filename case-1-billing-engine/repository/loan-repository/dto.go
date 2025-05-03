package loanrepository

import (
	"database/sql"
	"time"

	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/entity"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Loan struct {
	ID                int64          `gorm:"column:id;primary_key"`
	Serial            string         `gorm:"column:serial"`
	CreatedBy         string         `gorm:"column:created_by"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedBy         string         `gorm:"column:updated_by"`
	UpdatedAt         time.Time      `gorm:"column:updated_at"`
	DeletedBy         sql.NullString `gorm:"column:deleted_by"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at"`
	BorrowerName      string         `gorm:"column:borrower_name"`
	LoanName          string         `gorm:"column:loan_name"`
	Description       string         `gorm:"column:description"`
	LoanAmount        float64        `gorm:"column:loan_amount"`
	InterestRate      float32        `gorm:"column:interest_rate"`
	LoanTerm          int32          `gorm:"column:loan_term"`
	LoanTermUnit      string         `gorm:"column:loan_term_unit"`
	LoanTotalAmount   float64        `gorm:"column:loan_total_amount"`
	LoanOutstanding   float64        `gorm:"column:loan_outstanding"`
	LoanPrincipalOnly float64        `gorm:"column:loan_principal_only"`
}

func (l *Loan) ToEntity() entity.Loan {
	return entity.Loan{
		ID:                l.ID,
		Serial:            l.Serial,
		BorrowerName:      l.BorrowerName,
		LoanName:          l.LoanName,
		Description:       l.Description,
		LoanAmount:        l.LoanAmount,
		InterestRate:      l.InterestRate,
		LoanTerm:          l.LoanTerm,
		LoanTermUnit:      l.LoanTermUnit,
		LoanTotalAmount:   l.LoanTotalAmount,
		LoanOutstanding:   l.LoanOutstanding,
		LoanPrincipalOnly: l.LoanPrincipalOnly,
	}
}

func (l *Loan) FromEntity(record entity.Loan) {
	l.ID = record.ID
	l.Serial = record.Serial
	l.BorrowerName = record.BorrowerName
	l.LoanName = record.LoanName
	l.Description = record.Description
	l.LoanAmount = record.LoanAmount
	l.InterestRate = record.InterestRate
	l.LoanTerm = record.LoanTerm
	l.LoanTermUnit = record.LoanTermUnit
	l.LoanTotalAmount = record.LoanTotalAmount
	l.LoanOutstanding = record.LoanOutstanding
	l.LoanPrincipalOnly = record.LoanPrincipalOnly
}

type LoanRepayment struct {
	ID                         int64          `gorm:"column:id;primary_key"`
	Serial                     string         `gorm:"column:serial"`
	CreatedBy                  string         `gorm:"column:created_by"`
	CreatedAt                  time.Time      `gorm:"column:created_at"`
	UpdatedBy                  string         `gorm:"column:updated_by"`
	UpdatedAt                  time.Time      `gorm:"column:updated_at"`
	DeletedBy                  sql.NullString `gorm:"column:deleted_by"`
	DeletedAt                  gorm.DeletedAt `gorm:"column:deleted_at"`
	LoanSerial                 string         `gorm:"column:loan_serial"`
	InstallmentNumber          int32          `gorm:"column:installment_number"`
	Amount                     float64        `gorm:"column:amount"`
	InterestAmount             float64        `gorm:"column:interest_amount"`
	TotalAmount                float64        `gorm:"column:total_amount"`
	Status                     string         `gorm:"column:status"`
	RepaymentAmount            float64        `gorm:"column:repayment_amount"`
	RepaymentTime              time.Time      `gorm:"column:repayment_time"`
	RepaymentDueDate           time.Time      `gorm:"column:repayment_due_date"`
	LoanRepaymentHistorySerial string         `gorm:"column:loan_repayment_history_serial"`
}

func (lr *LoanRepayment) ToEntity() entity.LoanRepayment {
	return entity.LoanRepayment{
		ID:                         lr.ID,
		Serial:                     lr.Serial,
		Loan:                       entity.Loan{Serial: lr.LoanSerial},
		InstallmentNumber:          lr.InstallmentNumber,
		Amount:                     lr.Amount,
		InterestAmount:             lr.InterestAmount,
		TotalAmount:                lr.TotalAmount,
		Status:                     lr.Status,
		RepaymentAmount:            lr.RepaymentAmount,
		RepaymentTime:              lr.RepaymentTime,
		RepaymentDueDate:           lr.RepaymentDueDate,
		LoanRepaymentHistorySerial: lr.LoanRepaymentHistorySerial,
	}
}

func (lr *LoanRepayment) FromEntity(record entity.LoanRepayment) {
	lr.ID = record.ID
	lr.Serial = record.Serial
	lr.LoanSerial = record.Loan.Serial
	lr.InstallmentNumber = record.InstallmentNumber
	lr.Amount = record.Amount
	lr.InterestAmount = record.InterestAmount
	lr.TotalAmount = record.TotalAmount
	lr.Status = record.Status
	lr.RepaymentAmount = record.RepaymentAmount
	lr.RepaymentTime = record.RepaymentTime
	lr.RepaymentDueDate = record.RepaymentDueDate
	lr.LoanRepaymentHistorySerial = record.LoanRepaymentHistorySerial
}

type LoanRepaymentHistory struct {
	ID                          int64          `gorm:"column:id;primary_key"`
	Serial                      string         `gorm:"column:serial;type:uuid;default:uuid_generate_v4();uniqueIndex;<-:create"`
	CreatedBy                   string         `gorm:"column:created_by"`
	CreatedAt                   time.Time      `gorm:"column:created_at"`
	UpdatedBy                   string         `gorm:"column:updated_by"`
	UpdatedAt                   time.Time      `gorm:"column:updated_at"`
	DeletedBy                   sql.NullString `gorm:"column:deleted_by"`
	DeletedAt                   gorm.DeletedAt `gorm:"column:deleted_at"`
	LoanSerial                  string         `gorm:"column:loan_serial"`
	LoanRepaymentSerial         string         `gorm:"column:loan_repayment_serial"`
	RepaymentAmount             float64        `gorm:"column:repayment_amount"`
	RepaymentTime               time.Time      `gorm:"column:repayment_time"`
	RelatedLoanRepaymentSerials pq.StringArray `gorm:"type:uuid[];column:related_loan_repayment_serials"`
	Status                      string         `gorm:"column:status"`
}

func (lrh *LoanRepaymentHistory) ToEntity() entity.LoanRepaymentHistory {
	entityRecord := entity.LoanRepaymentHistory{
		ID:                    lrh.ID,
		Serial:                lrh.Serial,
		Loan:                  entity.Loan{Serial: lrh.LoanSerial},
		LoanRepayment:         entity.LoanRepayment{Serial: lrh.LoanRepaymentSerial},
		RepaymentAmount:       lrh.RepaymentAmount,
		RepaymentTime:         lrh.RepaymentTime,
		RelatedLoanRepayments: []entity.LoanRepayment{},
		Status:                lrh.Status,
	}

	for _, serial := range lrh.RelatedLoanRepaymentSerials {
		entityRecord.RelatedLoanRepayments = append(entityRecord.RelatedLoanRepayments, entity.LoanRepayment{Serial: serial})
	}

	return entityRecord
}

func (lrh *LoanRepaymentHistory) FromEntity(record entity.LoanRepaymentHistory) {
	lrh.ID = record.ID
	lrh.Serial = record.Serial
	lrh.LoanSerial = record.Loan.Serial
	lrh.LoanRepaymentSerial = record.LoanRepayment.Serial
	lrh.RepaymentAmount = record.RepaymentAmount
	lrh.RepaymentTime = record.RepaymentTime
	lrh.Status = record.Status

	lrh.RelatedLoanRepaymentSerials = pq.StringArray{}
	for _, relatedLoanRepayment := range record.RelatedLoanRepayments {
		lrh.RelatedLoanRepaymentSerials = append(lrh.RelatedLoanRepaymentSerials, relatedLoanRepayment.Serial)
	}
}
