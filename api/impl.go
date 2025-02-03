package api

import (
	"log-backend/internal/db"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func NewServer() Server {

	db, err := db.InitDB()
	if err != nil {
		logrus.Fatalf("init db error: %v", err)
		panic(err)
	}

	return Server{db}
}
