package entity

import (
	"time"
)

const (
	LoanRepaymentStatusPending  = "pending"
	LoanRepaymentStatusPaid     = "paid"
	LoanRepaymentStatusUnpaid   = "unpaid"
	LoanRepaymentStatusCanceled = "canceled"
)

type Loan struct {
	ID                int64   `json:"id"`
	Serial            string  `json:"serial"`
	BorrowerName      string  `json:"borrower_name"`
	LoanName          string  `json:"loan_name"`
	Description       string  `json:"description"`
	LoanAmount        float64 `json:"loan_amount"`
	InterestRate      float32 `json:"interest_rate"`
	LoanTerm          int32   `json:"loan_term"`
	LoanTermUnit      string  `json:"loan_term_unit"`
	LoanTotalAmount   float64 `json:"loan_total_amount"`
	LoanOutstanding   float64 `json:"loan_outstanding"`
	LoanPrincipalOnly float64 `json:"loan_principal_only"`
}

type LoanRepayment struct {
	ID                         int64     `json:"id"`
	Serial                     string    `json:"serial"`
	Loan                       Loan      `json:"loan"`
	InstallmentNumber          int32     `json:"installment_number"`
	Amount                     float64   `json:"amount"`
	InterestAmount             float64   `json:"interest_amount"`
	TotalAmount                float64   `json:"total_amount"`
	Status                     string    `json:"status"`
	RepaymentAmount            float64   `json:"repayment_amount"`
	RepaymentTime              time.Time `json:"repayment_time"`
	RepaymentDueDate           time.Time `json:"repayment_due_date"`
	LoanRepaymentHistorySerial string    `json:"loan_repayment_history_serial"`
}

type LoanRepaymentHistory struct {
	ID                    int64           `json:"id"`
	Serial                string          `json:"serial"`
	Loan                  Loan            `json:"loan"`
	LoanRepayment         LoanRepayment   `json:"loan_repayment"`
	RepaymentAmount       float64         `json:"repayment_amount"`
	RepaymentTime         time.Time       `json:"repayment_time"`
	RelatedLoanRepayments []LoanRepayment `json:"related_loan_repayments"`
	Status                string          `json:"status"`
}

type UpsertRepaymentTransactionRequest struct {
	Loan                 Loan                 `json:"loan"`
	LoanRepaymentHistory LoanRepaymentHistory `json:"loan_repayment_history"`
	LoanRepayments       []LoanRepayment      `json:"loan_repayment"`
}

type MakeRepaymentRequest struct {
	LoanSerial string    `json:"loan_serial"`
	DateStr    string    `json:"date"`
	Date       time.Time `json:"-"`
	Status     string    `json:"status"`
}

type MakeRepaymentResponse struct {
	RepaymentTransactionStatus bool            `json:"repayment_status"`
	Message                    string          `json:"message"`
	LoanRepayments             []LoanRepayment `json:"loan_repayment"`
}

type GetLoanRepaymentsRequest struct {
	LoanSerial string    `json:"loan_serial"`
	MinDate    time.Time `json:"min_date"`
	MaxDate    time.Time `json:"max_date"`
	Serials    []string  `json:"serials"`
	Statuses   []string  `json:"statuses"`
}

type IsDelinquentByLoanSerialRequest struct {
	LoanSerial string    `json:"loan_serial"`
	MaxDateStr string    `json:"max_date"` // ubah jadi string sementara
	MaxDate    time.Time `json:"-"`
}

type IsDelinquentByLoanSerialResponse struct {
	IsDelinquent         bool            `json:"is_delinquent"`
	TotalUnpaidAmount    float64         `json:"total_unpaid_amount"`
	UnpaidLoanRepayments []LoanRepayment `json:"unpaid_loan_repayments"`
}

type GetOutstandingByLoanSerialRequest struct {
	LoanSerial string `json:"loan_serial"`
}

type GetLoanRepaymentHistoryRequest struct {
	Serial     string `json:"serial"`
	LoanSerial string `json:"loan_serial"`
	Status     string `json:"status"`
}

type GetRepaymentListRequest struct {
	Page         int32     `json:"page"`
	PageSize     int32     `json:"page_size"`
	LoanSerial   string    `json:"loan_serial"`
	StartDateStr string    `json:"start_date"`
	StartDate    time.Time `json:"-"`
	EndDateStr   string    `json:"end_date_date"`
	EndDate      time.Time `json:"-"`
	Statuses     []string  `json:"statuses"`
}

type GetRepaymentListResponse struct {
	Page      int32           `json:"page"`
	PageSize  int32           `json:"page_size"`
	TotalData int64           `json:"total_data"`
	TotalPage int64           `json:"total_page"`
	Items     []LoanRepayment `json:"items"`
}

type GetUnpaidLoanRepaymentRequest struct {
	DateStr    string    `json:"date"`
	Date       time.Time `json:"-"`
	LoanSerial string    `json:"loan_serial"`
}

type GetRepaymentHistoryListRequest struct {
	Page         int32     `json:"page"`
	PageSize     int32     `json:"page_size"`
	LoanSerial   string    `json:"loan_serial"`
	StartDateStr string    `json:"start_date"`
	StartDate    time.Time `json:"-"`
	EndDateStr   string    `json:"end_date_date"`
	EndDate      time.Time `json:"-"`
	Statuses     []string  `json:"statuses"`
}

type GetRepaymentHistoryListResponse struct {
	Page      int32                  `json:"page"`
	PageSize  int32                  `json:"page_size"`
	TotalData int64                  `json:"total_data"`
	TotalPage int64                  `json:"total_page"`
	Items     []LoanRepaymentHistory `json:"items"`
}
