package main

import (
	"context"
	"fmt"
	"log-backend/api"
	"log-backend/autogen"
	"log-backend/internal/db"
	"os"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oapi-codegen/echo-middleware"
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
	server.Init()

	e := echo.New()
	e.Debug = true

	autogen.RegisterHandlers(e, server)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:5173", "https://log.ceten.fr"}, // Remplace "*" par une origine spécifique si nécessaire
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
	}))

	// Global middlewares
	swagger, _ := autogen.GetSwagger()
	e.Use(echomiddleware.OapiRequestValidatorWithOptions(swagger, &echomiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				// Convert the echo.Context from the context
				echoCtx := echomiddleware.GetEchoContext(ctx)
				if echoCtx == nil {
					return fmt.Errorf("error getting echo context")
				}
				
				// Call our handler
				return server.SecurityHandler(echoCtx, input.SecuritySchemeName, input.Scopes)
			},
		},
	}))

	// And we serve HTTP until the world ends.
	logrus.Fatal(e.Start("0.0.0.0:" + os.Getenv("PORT")))
}
