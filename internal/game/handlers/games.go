package handlers

import (
	"juanmagc99/checkers/internal/game/models"
	"juanmagc99/checkers/internal/game/utils"
	"juanmagc99/checkers/internal/storage"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GameHandler struct {
	redisStore storage.RedisStore
}

func NewGameHandler(redisStore storage.RedisStore) *GameHandler {
	return &GameHandler{
		redisStore: redisStore,
	}
}

func (h *GameHandler) CreateGame(c echo.Context) error {

	gameID := uuid.New().String()

	game := models.Game{
		ID:     gameID,
		Board:  models.InitBoard(),
		Turn:   "W",
		Status: "waiting",
	}

	if err := h.redisStore.Set(c.Request().Context(), gameID, game, 1*time.Hour); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Error saving game in database")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "Success initializing a new game",
		"game_id": gameID,
	})
}
