package db

import (
	"database/sql"
	"errors"
	"log-backend/internal/models"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*sql.DB, *gorm.DB, error) {
    dsn := os.Getenv("LOG_BACKEND_DSN")
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        logrus.Fatalf("open db error: %v", err)
    }
    // Get the underlying sql.DB object to close the connection later
    sqlDB, err := db.DB()
    if err != nil {
        logrus.Printf("get db error: %v", err)
        return nil, nil, errors.New("get db error")
    }
    // Ping the database to check if the connection is successful
    err = sqlDB.Ping()
    if err != nil {
        logrus.Printf("ping db error: %v", err)
        return nil, nil, errors.New("ping db error")
    }
    logrus.Println("Database connection successful")
    // Perform auto migration
    err = db.AutoMigrate(models.Health{}, models.Account{})
    if err != nil {
        logrus.Printf("auto migrate error: %v", err)
        return nil, nil, errors.New("auto migrate error")
    }
    logrus.Println("Auto migration completed")

	health := models.Health{}
	result := db.First(&health)
	if result.Error == gorm.ErrRecordNotFound {
		health.Nb = 0
		health.Status = "OK"
		result_create := db.Create(&health)
		if result_create.Error != nil {
			logrus.Fatalf("create health error: %v", result_create.Error)
			return nil, nil, errors.New("create health error")
		}
	}
	return sqlDB, db, nil
}