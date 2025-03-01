package main

import (
	"context"
	"juanmagc99/checkers/internal/game/handlers"
	"juanmagc99/checkers/internal/game/routes"
	"juanmagc99/checkers/internal/storage"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	e := echo.New()

	e.Logger.SetOutput(os.Stdout)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	redisStore := storage.NewRedisStore(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := redisStore.Set(ctx, "test_key", "test_value", 10*time.Second); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	gameHandler := handlers.NewGameHandler(redisStore)

	routes.RegisterRoutes(e, gameHandler)

	err := e.Start("localhost:8080")

	if err != nil {
		e.Logger.Fatal("Error starting server")
	}
}
