package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"log-backend/autogen"
)

// optional code omitted

type Server struct{}

func NewServer() Server {
	return Server{}
}

// (GET /ping)
func (Server) GetPing(ctx echo.Context) error {
	resp := autogen.Pong{
		Ping: "pong",
	}

	return ctx.JSON(http.StatusOK, resp)
}