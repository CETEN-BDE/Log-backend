package api

import (

	"log-backend/autogen"
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
		autogen.GetHealth500JSONResponse{Message: "Error getting health"}.VisitGetHealthResponse(ctx.Response())
		return result.Error
	}

	health.Nb += 1
	result = s.db.Save(&health)
	if result.Error != nil {
		logrus.Errorf("Error updating health: %v", result.Error)
		autogen.GetHealth500JSONResponse{Message: "Error updating health"}.VisitGetHealthResponse(ctx.Response())
		return result.Error
	}

	logrus.Info("/health: Health check")
	autogen.GetHealth200JSONResponse{Status: health.Status, Nb: health.Nb}.VisitGetHealthResponse(ctx.Response())
	return nil
}
