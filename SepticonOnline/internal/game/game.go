package game

import (
	"septikon-online/internal/player"
)

// Статус игровой партии
type GameStatus string

const (
	StatusWaiting  GameStatus = "waiting"
	StatusPlaying  GameStatus = "playing"
	StatusFinished GameStatus = "finished"
)

// Информация о текущем ходе
type Turn struct {
	CurrentPlayerID string `json:"currentPlayerId"`
	Number          int    `json:"number"`
	DiceResult      int    `json:"diceResult"`
	ActionUsed      bool   `json:"actionUsed"`
}

// Основная структура партии
type Game struct {
	ID       string          `json:"id"`
	Status   GameStatus      `json:"status"`
	Players  []player.Player `json:"players"`
	Turn     Turn            `json:"turn"`
	WinnerID string          `json:"winnerId,omitempty"`
}
