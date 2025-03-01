package routes

import (
	"juanmagc99/checkers/internal/game/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo,
	gameHandler *handlers.GameHandler) {

	api := e.Group("/api")

	api.POST("/games", gameHandler.CreateGame)
	api.POST("/games/:id/join", gameHandler.JoinGame)
	api.GET("/games/:id", gameHandler.GetGame)
}
