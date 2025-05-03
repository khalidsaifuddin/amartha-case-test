package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/core/module"
	_ "github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/docs"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/handler/api"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/handler/worker"
	workerclient "github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/handler/worker/client"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/pkg/conn"
	transactionrepository "github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/repository/transaction-repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(cfg config.Config, db *gorm.DB) (*gin.Engine, conn.CacheService) {
	if strings.EqualFold(cfg.Environment, "production") {
		gin.SetMode(gin.ReleaseMode)
	}

	coreRedis, redisPool := conn.InitRedis(cfg)
	redisClientWorker := worker.InitRedisClientWorker(cfg)

	// initiate worker client
	workerClient := workerclient.NewWorkerClient(asynq.NewClient(redisClientWorker))

	// initialize repository
	transactionRepo := transactionrepository.New(cfg, redisPool)

	// initialize usecase
	transactionUc := module.NewTransactionUsecase(cfg, transactionRepo, workerClient)

	// initialize worker server
	workerServer := worker.NewWorkerServer(redisClientWorker, cfg, transactionUc)

	// initiaalize http handler
	httpHandler := api.NewHTTPHandler(cfg, transactionUc)

	if cfg.WorkerEnabled {
		go workerServer.Run()
	}

	if cfg.SchedulerEnabled {
		go workerServer.RunScheduler()
	}

	router := gin.New()
	router.Use(CORSMiddleware(cfg))

	// router config
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// router group for v1
	v1 := router.Group("/v1")
	{
		v1.POST("/transaction/reconcile", httpHandler.ReconcileTransaction)
		v1.GET("/transaction/reconcile/:serial", httpHandler.GetTransactionReconciliationProgress)
	}

	// API documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "Page not found"})
	})

	return router, coreRedis
}
