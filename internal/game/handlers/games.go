package handlers

import (
	"juanmagc99/checkers/internal/game/models"
	"juanmagc99/checkers/internal/game/utils"
	"juanmagc99/checkers/internal/storage"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

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

func (h *GameHandler) GameSessionHandler(c echo.Context) error {

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

	pRole, err := game.ToSafeGame(playerToken)
	if err != nil {
		c.Logger().Error(err)
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized access")
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error(err)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
	}

	gameSession := models.GetSession(gameID)

	gameSession.Mu.Lock()
	defer gameSession.Mu.Unlock()

	if pRole == "player1" {
		if gameSession.Player1Conn != nil {
			gameSession.Player1Conn.Close()
		}
		gameSession.Player1Conn = conn
	} else if pRole == "player2" {
		if gameSession.Player2Conn != nil {
			gameSession.Player2Conn.Close()
		}
		gameSession.Player2Conn = conn
	} else {
		conn.Close()
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid player token")
	}

	go forwardMessages(gameSession, pRole)

	select {}
}

func forwardMessages(gs *models.GameSession, playerRole string) {
	// Log inicial: indicar que se inició el reenvío para este rol.
	log.Printf("[forwardMessages] Iniciando reenvío para el rol: %s", playerRole)

	for {
		var src, dest *websocket.Conn

		gs.Mu.Lock()
		if playerRole == "player1" {
			src = gs.Player1Conn
			dest = gs.Player2Conn
			log.Printf("[forwardMessages] (player1) src asignado, dest asignado.")
		} else {
			src = gs.Player2Conn
			dest = gs.Player1Conn
			log.Printf("[forwardMessages] (player2) src asignado, dest asignado.")
		}
		gs.Mu.Unlock()

		if src == nil || dest == nil {
			log.Printf("[forwardMessages] Una de las conexiones es nil (src: %v, dest: %v). Esperando...", src, dest)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		log.Printf("[forwardMessages] Esperando mensaje de %s...", playerRole)
		_, msg, err := src.ReadMessage()
		if err != nil {
			log.Printf("[forwardMessages] Error al leer mensaje desde %s: %v", playerRole, err)
			src.Close()
			return
		}
		log.Printf("[forwardMessages] Mensaje recibido de %s: %s", playerRole, string(msg))

		err = dest.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("[forwardMessages] Error al escribir mensaje al destino: %v", err)
			dest.Close()
			return
		}
		log.Printf("[forwardMessages] Mensaje reenviado correctamente.")
	}
}
