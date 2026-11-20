package game

import (
	"errors"

	"septikon-online/internal/player"
	"septikon-online/internal/rocket"
	"septikon-online/internal/turns"
)

// Статус игровой партии
type GameStatus string

const (
	StatusWaiting  GameStatus = "waiting"
	StatusPlaying  GameStatus = "playing"
	StatusFinished GameStatus = "finished"
)

// Основная структура партии
type Game struct {
	ID          string          `json:"id"`
	Status      GameStatus      `json:"status"`
	Players     []player.Player `json:"players"`
	Rockets     []rocket.Rocket `json:"rockets"`
	TurnManager *turns.Manager  `json:"turn"`
	WinnerID    string          `json:"winnerId,omitempty"`
	Repairs     []RepairTask    `json:"repairs"`
}

// Создание новой партии
func NewGame(
	gameID string,
	player1ID string,
	player1Name string,
	player2ID string,
	player2Name string,
) (*Game, error) {

	if gameID == "" || player1ID == "" || player2ID == "" {
		return nil, errors.New("ID игры и игроков не могут быть пустыми")
	}

	if player1ID == player2ID {
		return nil, errors.New("игроки должны иметь разные ID")
	}

	if player1Name == "" || player2Name == "" {
		return nil, errors.New("имена игроков не могут быть пустыми")
	}

	players := []player.Player{
		player.NewPlayer(player1ID, player1Name),
		player.NewPlayer(player2ID, player2Name),
	}

	turnManager, err := turns.NewManager(
		[]string{player1ID, player2ID},
	)

	if err != nil {
		return nil, err
	}

	game := &Game{
		ID:          gameID,
		Status:      StatusPlaying,
		Players:     players,
		Rockets:     make([]rocket.Rocket, 0),
		TurnManager: turnManager,
		Repairs:     make([]RepairTask, 0),
	}

	return game, nil
}

// Поиск игрока по ID
func (g *Game) FindPlayer(playerID string) (*player.Player, error) {
	for i := range g.Players {
		if g.Players[i].ID == playerID {
			return &g.Players[i], nil
		}
	}

	return nil, errors.New("игрок не найден")
}

// Проверка текущего хода
func (g *Game) CheckTurn(playerID string) error {
	if g.Status != StatusPlaying {
		return errors.New("партия не активна")
	}

	if !g.TurnManager.IsPlayerTurn(playerID) {
		return errors.New("сейчас ход другого игрока")
	}

	return nil
}
