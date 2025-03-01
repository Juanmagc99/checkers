package utils

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetGameID(c echo.Context) (string, error) {
	gameID := c.Param("id")
	if gameID == "" {
		return "", errors.New("game id not provided")
	}

	if _, err := uuid.Parse(gameID); err != nil {
		return "", errors.New("invalid game id format")
	}
	return gameID, nil
}
