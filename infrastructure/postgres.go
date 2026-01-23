package infrastructure

import (
	"iam-service/config"
	"log"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg config.PostgresConfig, logger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.Platform.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Postgres")

	return db, nil
}
