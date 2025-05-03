package conn

import (
	"fmt"
	"log"
	"time"

	"github.com/khalidsaifuddin/amartha-case-test/case-1-billing-engine/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta", cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port)
	log.Printf("Attempting database connection using credentials: %v", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("%s", err.Error())
		panic(err)
	} else {
		log.Printf("Successfully connected to database server")
	}

	rdb, err := db.DB()
	if err != nil {
		log.Fatalf("%s", err.Error())
		panic(err)
	}

	rdb.SetMaxIdleConns(cfg.MaxIdleConns)
	rdb.SetMaxOpenConns(cfg.MaxOpenConns)
	rdb.SetConnMaxLifetime(time.Duration(int(time.Minute) * cfg.ConnMaxLifetime))

	return db
}

func DbClose(db *gorm.DB) {
	rdb, err := db.DB()
	if err != nil {
		log.Fatalf("%s", err.Error())
		panic(err)
	}

	_ = rdb.Close()
}
