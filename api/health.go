package api

import (
	"net/http"

	"log-backend/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// (GET /health)
func (s Server) GetHealth(ctx echo.Context) error {

	var health models.Health

	result := s.db.First(&health)
	if result.Error != nil {
		logrus.Errorf("Error getting health: %v", result.Error)
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}

	health.Nb += 1
	result = s.db.Save(&health)
	if result.Error != nil {
		logrus.Errorf("Error updating health: %v", result.Error)
		return ctx.JSON(http.StatusInternalServerError, result.Error)
	}

	logrus.Info("/health: Health check")
	return ctx.JSON(http.StatusOK, health)
}
