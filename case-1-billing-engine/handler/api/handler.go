package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/entity"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/module"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/pkg/helper"
)

type HTTPHandler interface {
	MakeRepayment(c *gin.Context)
	UpsertRepaymentTransaction(c *gin.Context)
	IsDelinquent(c *gin.Context)
	GetOutstandingByLoanSerial(c *gin.Context)
	GetRepaymentList(c *gin.Context)
	GetNextRepayment(c *gin.Context)
	GetRepaymentHistoryList(c *gin.Context)
}

type httpHandler struct {
	cfg    config.Config
	loanUc module.LoanUsecase
}

func NewHTTPHandler(cfg config.Config, loanUc module.LoanUsecase) HTTPHandler {
	return &httpHandler{
		cfg:    cfg,
		loanUc: loanUc,
	}
}

// MakeRepayment godoc
// @Summary Make a loan repayment
// @Description Apply a repayment to the loan with the given serial
// @Tags Loan
// @Accept json
// @Produce json
// @Param loan_serial path string true "Loan Serial"
// @Param request body entity.MakeRepaymentRequest true "Make Repayment payload"
// @Success 200 {object} entity.MakeRepaymentResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1/loan/repayment/{loan_serial} [post]
func (h *httpHandler) MakeRepayment(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	// bind request to struct
	request := entity.MakeRepaymentRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	parsedDate, err := time.Parse(entity.DateFormat, request.DateStr)
	if err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = "Invalid date format. Use YYYY-MM-DD"

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// inject parsed date into max date
	request.Date = parsedDate

	loanSerial := c.Param("loan_serial")
	if loanSerial == "" {
		statusCode = http.StatusBadRequest
		statusMessage = entity.ErrLoanSerialIsRequired.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	request.LoanSerial = loanSerial

	// call usecase
	resp, err := h.loanUc.MakeRepayment(c, request)
	if err != nil {
		statusCode = http.StatusInternalServerError
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// response
	helper.ResponseOutput(c, statusCode, statusMessage, resp)
}

// UpsertRepaymentTransaction godoc
// @Summary Manual transaction for repayment process
// @Description Provide manual transaction for repayment process for FE if custom payload is enabled
// @Tags Loan
// @Accept json
// @Produce json
// @Param loan_serial path string true "Loan Serial"
// @Param request body entity.UpsertRepaymentTransactionRequest true "UpsertRepaymentTransaction payload"
// @Success 200 {object} entity.UpsertRepaymentTransactionRequest
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1/loan/repayment/{loan_serial}/manual [post]
func (h *httpHandler) UpsertRepaymentTransaction(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	// bind request to struct
	request := entity.UpsertRepaymentTransactionRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	loanSerial := c.Param("loan_serial")
	if loanSerial == "" {
		statusCode = http.StatusBadRequest
		statusMessage = entity.ErrLoanSerialIsRequired.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	request.Loan.Serial = loanSerial

	// call usecase
	resp, err := h.loanUc.UpsertRepaymentTransaction(c, request)
	if err != nil {
		statusCode = http.StatusInternalServerError
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// response
	helper.ResponseOutput(c, statusCode, statusMessage, resp)
}

// IsDelinquent godoc
// @Summary Check if Loan Borrower is specified as delinquent
// @Description Return true of there are at least 2 unpaid repayment
// @Tags Loan
// @Accept json
// @Produce json
// @Param loan_serial path string true "Loan Serial"
// @Param request body entity.IsDelinquentByLoanSerialRequest true "Is Delinquent payload"
// @Success 200 {object} entity.IsDelinquentByLoanSerialResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1/loan/is_delinquent/{loan_serial} [post]
func (h *httpHandler) IsDelinquent(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	request := entity.IsDelinquentByLoanSerialRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	parsedDate, err := time.Parse(entity.DateFormat, request.MaxDateStr)
	if err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = "Invalid date format. Use YYYY-MM-DD"

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// inject parsed date into max date
	request.MaxDate = parsedDate

	// enforce loan serial from url param
	if c.Param("loan_serial") != "" {
		request.LoanSerial = c.Param("loan_serial")
	}

	// call usecase
	resp, err := h.loanUc.IsDelinquentByLoanSerial(c, request)
	if err != nil {
		statusCode = http.StatusInternalServerError
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// response
	helper.ResponseOutput(c, statusCode, statusMessage, resp)
}

// GetOutstandingByLoanSerial godoc
// @Summary Get Loan Outstanding
// @Description Get Loan Outstanding
// @Tags Loan
// @Accept json
// @Produce json
// @Param loan_serial path string true "Loan Serial"
// @Param request body entity.GetOutstandingByLoanSerialRequest true "GetOutstandingByLoanSerial payload"
// @Success 200 {object} entity.Loan
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1/loan/outstanding/{loan_serial} [get]
func (h *httpHandler) GetOutstandingByLoanSerial(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	request := entity.GetOutstandingByLoanSerialRequest{}

	// enforce loan serial from url param
	if c.Param("loan_serial") == "" {
		statusCode = http.StatusBadRequest
		statusMessage = entity.ErrLoanSerialIsRequired.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	request.LoanSerial = c.Param("loan_serial")

	// call usecase
	resp, err := h.loanUc.GetOutstandingByLoanSerial(c, request)
	if err != nil {
		statusCode = http.StatusInternalServerError
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// response
	helper.ResponseOutput(c, statusCode, statusMessage, resp)
}

// GetRepaymentList godoc
// @Summary Get List of repayment schedule for one loan
// @Description Get List of repayment schedule for one loan
// @Tags Loan
// @Accept json
// @Produce json
// @Param loan_serial path string true "Loan Serial"
// @Success 200 {object} entity.GetRepaymentListResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1/loan/repayment/{loan_serial} [get]
func (h *httpHandler) GetRepaymentList(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	request := entity.GetRepaymentListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// enforce loan serial from url param
	if c.Param("loan_serial") == "" {
		statusCode = http.StatusBadRequest
		statusMessage = entity.ErrLoanSerialIsRequired.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	request.LoanSerial = c.Param("loan_serial")

	// parse start date
	if request.StartDateStr != "" {
		parsedStartDate, err := time.Parse(entity.DateFormat, request.StartDateStr)
		if err != nil {
			statusCode = http.StatusBadRequest
			statusMessage = "Invalid date format. Use YYYY-MM-DD"

			log.Println(statusMessage)
			helper.ResponseOutput(c, statusCode, statusMessage, nil)
			return
		}

		request.StartDate = parsedStartDate
	}

	// parse end date
	if request.EndDateStr != "" {
		parsedEndDate, err := time.Parse(entity.DateFormat, request.EndDateStr)
		if err != nil {
			statusCode = http.StatusBadRequest
			statusMessage = "Invalid date format. Use YYYY-MM-DD"

			log.Println(statusMessage)
			helper.ResponseOutput(c, statusCode, statusMessage, nil)
			return
		}

		request.EndDate = parsedEndDate
	}

	// call usecase
	resp, err := h.loanUc.GetRepaymentList(c, request)
	if err != nil {
		statusCode = http.StatusInternalServerError
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// response
	helper.ResponseOutput(c, statusCode, statusMessage, resp)
}

// GetNextRepayment godoc
// @Summary Get List of next repayment schedule for one loan
// @Description Get List of next repayment schedule for one loan
// @Tags Loan
// @Accept json
// @Produce json
// @Param loan_serial path string true "Loan Serial"
// @Success 200 {object} []entity.LoanRepayment
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1/loan/repayment/next-schedule/{loan_serial} [get]
func (h *httpHandler) GetNextRepayment(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	request := entity.GetUnpaidLoanRepaymentRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// parse start date
	if request.DateStr != "" {
		parsedDate, err := time.Parse(entity.DateFormat, request.DateStr)
		if err != nil {
			statusCode = http.StatusBadRequest
			statusMessage = "Invalid date format. Use YYYY-MM-DD"

			log.Println(statusMessage)
			helper.ResponseOutput(c, statusCode, statusMessage, nil)
			return
		}

		request.Date = parsedDate
	}

	// call usecase
	resp, err := h.loanUc.GetNextRepayment(c, request)
	if err != nil {
		statusCode = http.StatusInternalServerError
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// response
	helper.ResponseOutput(c, statusCode, statusMessage, resp)
}

// GetRepaymentHistoryList godoc
// @Summary Get List of repayment history for one loan
// @Description Get List of repayment history for one loan
// @Tags Loan
// @Accept json
// @Produce json
// @Param loan_serial path string true "Loan Serial"
// @Success 200 {object} entity.GetRepaymentHistoryListResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1/loan/repayment-history/{loan_serial} [get]
func (h *httpHandler) GetRepaymentHistoryList(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	request := entity.GetRepaymentHistoryListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// parse start date
	if request.StartDateStr != "" {
		parsedStartDate, err := time.Parse(entity.DateFormat, request.StartDateStr)
		if err != nil {
			statusCode = http.StatusBadRequest
			statusMessage = "Invalid date format. Use YYYY-MM-DD"

			log.Println(statusMessage)
			helper.ResponseOutput(c, statusCode, statusMessage, nil)
			return
		}

		request.StartDate = parsedStartDate
	}

	// parse end date
	if request.EndDateStr != "" {
		parsedEndDate, err := time.Parse(entity.DateFormat, request.EndDateStr)
		if err != nil {
			statusCode = http.StatusBadRequest
			statusMessage = "Invalid date format. Use YYYY-MM-DD"

			log.Println(statusMessage)
			helper.ResponseOutput(c, statusCode, statusMessage, nil)
			return
		}

		request.EndDate = parsedEndDate
	}

	// call usecase
	resp, err := h.loanUc.GetRepaymentHistoryList(c, request)
	if err != nil {
		statusCode = http.StatusInternalServerError
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// response
	helper.ResponseOutput(c, statusCode, statusMessage, resp)
}
