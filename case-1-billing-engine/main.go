package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/config"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/handler/middleware"
	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/pkg/conn"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system environment")
	}

	cfg := config.Get()

	fmt.Printf("starting %s Service...", cfg.ServiceName)

	db := conn.InitDB(&cfg)
	defer conn.DbClose(db)

	router, _ := middleware.InitRouter(cfg, db)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		panic(fmt.Errorf("failed to start server: %s", err.Error()))
	}
}
