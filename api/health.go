package api

import (
	"net/http"

	"log-backend/autogen"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// (GET /health)
func (Server) GetHealth(ctx echo.Context) error {
	resp := autogen.HealthCheck{
		Status: "ok",
	}
	
	logrus.Info("/health: Health check")
	return ctx.JSON(http.StatusOK, resp)
}