package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GameHandler struct {
}

func NewGameHandler() *GameHandler {
	return &GameHandler{}
}

func (h *GameHandler) CreateGame(c echo.Context) error {
	return c.JSON(http.StatusOK, "Game has been created here is the id ...")
}
