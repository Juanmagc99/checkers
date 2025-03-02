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

func (h *GameHandler) CreateGameHandler(c echo.Context) error {
	gameID := uuid.New().String()
	player1Token := uuid.New().String()

	game := models.Game{
		ID:           gameID,
		Board:        models.InitBoard(),
		Turn:         "W",
		Status:       models.StatusWaiting,
		Player1Token: player1Token,
		//Player2Token to be empty till someone joins the game
	}

	if err := h.redisStore.Set(c.Request().Context(), gameID, game, 1*time.Hour); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Error saving game in database")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message":       "Success initializing a new game",
		"game_id":       gameID,
		"player1_token": player1Token,
	})
}

func (h *GameHandler) JoinGameHandler(c echo.Context) error {
	gameID, err := utils.GetGameID(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	var game models.Game
	if err := h.redisStore.Get(ctx, gameID, &game); err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Game not found")
	}

	if game.Status != "waiting" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Game already in progress or finished")
	}

	player2Token := uuid.New().String()

	game.Player2Token = player2Token
	game.Status = "in_progress"
	game.Turn = "W"

	if err := h.redisStore.Set(ctx, gameID, game, 0); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update game")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":       "Successfully joined the game",
		"game_id":       gameID,
		"player2_token": player2Token,
	})
}

func (h *GameHandler) GetGameHandler(c echo.Context) error {
	gameID, err := utils.GetGameID(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	var game models.Game
	if err := h.redisStore.Get(ctx, gameID, &game); err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Game not found")
	}

	playerToken, err := utils.GetPlayerToken(c)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Player token not provided")
	}

	_, err = game.ToSafeGame(playerToken)
	if err != nil {
		c.Logger().Error(err)
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized access")
	}

	return c.JSON(http.StatusOK, game)
}
