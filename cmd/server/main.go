package main

import (
	"juanmagc99/checkers/internal/game/handlers"
	"juanmagc99/checkers/internal/game/routes"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Logger.SetOutput(os.Stdout)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	gameHandler := handlers.NewGameHandler()

	routes.RegisterRoutes(e, gameHandler)

	err := e.Start("localhost:8080")

	if err != nil {
		e.Logger.Fatal("Error starting server")
	}
}
