package database

import (
	"context"
	"log"
	"time"

	"github.com/jfernsio/slotswapper/internals/config"
	"github.com/jfernsio/slotswapper/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	config.LoadEnv()
	dsn := config.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("❌ failed to connect to db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ sql.DB error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("❌ db ping failed: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Event{}, &models.SwapRequest{}); err != nil {
		log.Fatalf("❌ migration failed: %v", err)
	}

	DB = db
	log.Println("✅ Database connected & migrated")
	//show databse tables and columns
	var tables []string
if err := db.Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
    panic(err)
}   
}
