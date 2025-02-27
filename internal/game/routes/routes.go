package routes

import (
	"juanmagc99/checkers/internal/game/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, gameHandler *handlers.GameHandler) {
	api := e.Group("/api")

	api.GET("/games", gameHandler.CreateGame)
}
