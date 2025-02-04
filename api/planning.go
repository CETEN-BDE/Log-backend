package api

import (
	"log-backend/autogen"

	"github.com/labstack/echo/v4"
)

func (s Server) GetPlanning(ctx echo.Context) error {
	autogen.GetPlanning200JSONResponse{Message: "Hello, World!"}.VisitGetPlanningResponse(ctx.Response())
	return nil
}