package entity

import "errors"

var (
	ErrLoanNotFound                = errors.New("loan is not found")
	ErrLoanSerialIsRequired        = errors.New("loan serial is required")
	ErrLoanRepaymentRecordRequired = errors.New("loan repayment record is required")
	ErrLoanReyamentAmountInvalid   = errors.New("there are invalid repayment amounts. details are attached to response")
	ErrNoUnpaidRepayment           = errors.New("no unpaid repayment is found for this loan")
)
