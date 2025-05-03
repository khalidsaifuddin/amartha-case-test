package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/core/module"
	_ "github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/docs"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/handler/api"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/pkg/conn"
	loanrepository "github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/repository/loan-repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(cfg config.Config, db *gorm.DB) (*gin.Engine, conn.CacheService) {
	if strings.EqualFold(cfg.Environment, "production") {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(CORSMiddleware(cfg))

	coreRedis, redis := conn.InitRedis(cfg)

	// repository initialization
	loanRepository := loanrepository.New(cfg, db, redis)

	// usecase initialization
	loanUsecase := module.NewLoanUsecase(cfg, loanRepository)

	// handler initialization
	httpHandler := api.NewHTTPHandler(cfg, loanUsecase)

	// router config
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// router group for v1
	v1 := router.Group("/v1")
	{
		v1.POST("/loan/repayment/manual/:loan-serial", httpHandler.UpsertRepaymentTransaction)
		v1.POST("/loan/repayment/:loan_serial", httpHandler.MakeRepayment)
		v1.POST("/loan/is_delinquent/:loan_serial", httpHandler.IsDelinquent)
		v1.GET("/loan/outstanding/:loan_serial", httpHandler.GetOutstandingByLoanSerial)
		v1.GET("/loan/repayment/:loan_serial", httpHandler.GetRepaymentList)
		v1.GET("/loan/repayment/next-schedule/:loan_serial", httpHandler.GetNextRepayment)
		v1.GET("/loan/repayment-history/:loan_serial", httpHandler.GetRepaymentHistoryList)
	}

	// API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "Page not found"})
	})

	return router, coreRedis
}
