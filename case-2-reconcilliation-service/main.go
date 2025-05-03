package main

import (
	"fmt"

	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/handler/middleware"
	"github.com/khalidsaifuddin/amartha-case-test/case-2-reconcilliation-service/pkg/conn"
)

func main() {
	cfg := config.Get()

	fmt.Printf("starting %s Service...", cfg.ServiceName)

	db := conn.InitDB(&cfg)
	defer conn.DbClose(db)

	router, _ := middleware.InitRouter(cfg, db)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		panic(fmt.Errorf("failed to start server: %s", err.Error()))
	}
}
