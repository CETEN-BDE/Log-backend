package main

import (
	"log-backend/api"
	"log-backend/autogen"
	"log-backend/internal/db"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	
	// Open a database connection
	sqlDB, db, err := db.InitDB()
	if err != nil {
		logrus.Fatalf("init db error: %v", err)
		panic(err)
	}

	defer sqlDB.Close()
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer(db)

	e := echo.New()

	autogen.RegisterHandlers(e, server)

	// And we serve HTTP until the world ends.
	logrus.Fatal(e.Start("0.0.0.0:" + os.Getenv("PORT")))
}
