package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/handler/worker"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/pkg/conn"
	"gorm.io/gorm"
)

func InitRouter(cfg config.Config, db *gorm.DB) (*gin.Engine, conn.CacheService) {
	if strings.EqualFold(cfg.Environment, "production") {
		gin.SetMode(gin.ReleaseMode)
	}

	redisClientWorker := worker.InitRedisClientWorker(cfg)
	workerServer := worker.NewWorkerServer(redisClientWorker, cfg)

	if cfg.WorkerEnabled {
		go workerServer.Run()
	}

	if cfg.SchedulerEnabled {
		go workerServer.RunScheduler()
	}

	router := gin.New()
	router.Use(CORSMiddleware(cfg))

	coreRedis, _ := conn.InitRedis(cfg)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "404", "message": "Page not found"})
	})

	return router, coreRedis
}
