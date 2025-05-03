package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/entity"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/module"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/pkg/helper"
)

type HTTPHandler interface {
	ReconcileTransaction(c *gin.Context)
	GetTransactionReconciliationProgress(c *gin.Context)
}

type httpHandler struct {
	cfg           config.Config
	transactionUc module.TransactionUsecase
}

func NewHTTPHandler(cfg config.Config, transactionUc module.TransactionUsecase) HTTPHandler {
	return &httpHandler{
		cfg:           cfg,
		transactionUc: transactionUc,
	}
}

// ReconcileTransaction godoc
// @Summary Reconcile Transaction between system and bank transaction
// @Description Reconcile Transaction between system and bank transaction
// @Tags Loan
// @Accept json
// @Produce json
// @Param request body entity.ReconcileTransactionRequest true "ReconcileTransaction payload"
// @Success 200 {object} entity.ReconcileTransactionResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /v1//transaction/reconcile [post]
func (h *httpHandler) ReconcileTransaction(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	// bind request to struct
	request := entity.ReconcileTransactionRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		statusCode = http.StatusBadRequest
		statusMessage = err.Error()

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	if request.StartDateStr != "" {
		parsedStartDate, err := time.Parse(entity.DateFormat, request.StartDateStr)
		if err != nil {
			statusCode = http.StatusBadRequest
			statusMessage = "Invalid date format. Use YYYY-MM-DD"

			log.Println(statusMessage)
			helper.ResponseOutput(c, statusCode, statusMessage, nil)
			return
		}

		// inject parsed date into date
		request.StartDate = parsedStartDate
	}

	if request.EndDateStr != "" {
		parsedEndDate, err := time.Parse(entity.DateFormat, request.EndDateStr)
		if err != nil {
			statusCode = http.StatusBadRequest
			statusMessage = "Invalid date format. Use YYYY-MM-DD"

			log.Println(statusMessage)
			helper.ResponseOutput(c, statusCode, statusMessage, nil)
			return
		}

		// inject parsed date into date
		request.EndDate = parsedEndDate
	}

	// call usecase
	resp, err := h.transactionUc.RunReconcileTransaction(c, request)
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

func (h *httpHandler) GetTransactionReconciliationProgress(c *gin.Context) {
	var statusCode int32 = http.StatusOK
	var statusMessage string = entity.DefaultSuccessMessage

	serial := c.Param("serial")
	if serial == "" {
		statusCode = http.StatusBadRequest
		statusMessage = "serial is required"

		log.Println(statusMessage)
		helper.ResponseOutput(c, statusCode, statusMessage, nil)
		return
	}

	// call usecase
	resp, err := h.transactionUc.GetTransactionReconciliationProgress(c, serial)
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
