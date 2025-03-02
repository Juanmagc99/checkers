package routes

import (
	"juanmagc99/checkers/internal/game/handlers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo,
	gameHandler *handlers.GameHandler) {

	api := e.Group("/api")

	api.POST("/games", gameHandler.CreateGameHandler)
	api.POST("/games/:id/join", gameHandler.JoinGameHandler)
	api.GET("/games/:id", gameHandler.GetGameHandler)
	api.GET("/games/:id/ws", gameHandler.GameSessionHandler)
}
