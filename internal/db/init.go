package db

import (
    "log"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
	"errors"
)

func InitDB() (*gorm.DB, error) {
    dsn := "root:changeme@tcp(127.0.0.1:3306)/log?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("open db error: %v", err)
    }
    // Get the underlying sql.DB object to close the connection later
    sqlDB, err := db.DB()
    if err != nil {
        log.Printf("get db error: %v", err)
        return nil, errors.New("get db error")
    }
    defer sqlDB.Close()
    // Ping the database to check if the connection is successful
    err = sqlDB.Ping()
    if err != nil {
        log.Printf("ping db error: %v", err)
        return nil, errors.New("ping db error")
    }
    log.Println("Database connection successful")
    // Perform auto migration
    err = db.AutoMigrate()
    if err != nil {
        log.Printf("auto migrate error: %v", err)
        return nil, errors.New("auto migrate error")
    }
    log.Println("Auto migration completed")
	return db, nil
}