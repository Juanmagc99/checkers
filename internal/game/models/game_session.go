package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type GameSession struct {
	Player1Conn *websocket.Conn
	Player2Conn *websocket.Conn
	Mu          sync.Mutex
}

var Sessions = make(map[string]*GameSession)
var sessionsMu sync.Mutex

func GetSession(gameID string) *GameSession {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	session, exists := Sessions[gameID]
	if !exists {
		session = &GameSession{}
		Sessions[gameID] = session
	}
	return session
}
