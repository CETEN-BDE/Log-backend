package api

import ( 
	"net/http"

	"github.com/labstack/echo/v4"
	"log-backend/autogen"
)

// (GET /health)
func (Server) GetHealth(ctx echo.Context) error {
	resp := autogen.HealthCheck{
		Status: "ok",
	}

	return ctx.JSON(http.StatusOK, resp)
}