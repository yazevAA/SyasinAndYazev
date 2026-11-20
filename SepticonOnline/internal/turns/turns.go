package turns

import (
	"errors"
	"math/rand"
	"time"
)

// Максимальное количество игроков
const MaxPlayers = 2

// Состояние игрового хода
type Turn struct {
	CurrentPlayerID string `json:"currentPlayerId"`
	Number          int    `json:"number"`
	DiceResult      int    `json:"diceResult"`
	DiceRolled      bool   `json:"diceRolled"`
	ActionUsed      bool   `json:"actionUsed"`
}

// Управление очередностью игроков
type Manager struct {
	PlayerIDs []string `json:"playerIds"`
	Turn      Turn     `json:"turn"`

	random *rand.Rand
}

// Создание системы ходов
func NewManager(playerIDs []string) (*Manager, error) {
	if len(playerIDs) != MaxPlayers {
		return nil, errors.New("для игры необходимо ровно два игрока")
	}

	if playerIDs[0] == "" || playerIDs[1] == "" {
		return nil, errors.New("ID игроков не могут быть пустыми")
	}

	if playerIDs[0] == playerIDs[1] {
		return nil, errors.New("игроки должны иметь разные ID")
	}

	manager := &Manager{
		PlayerIDs: playerIDs,
		Turn: Turn{
			CurrentPlayerID: playerIDs[0],
			Number:          1,
			DiceResult:      0,
			DiceRolled:      false,
			ActionUsed:      false,
		},
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return manager, nil
}

// Проверка, принадлежит ли ход игроку
func (m *Manager) IsPlayerTurn(playerID string) bool {
	return m.Turn.CurrentPlayerID == playerID
}

// Бросок кубика
func (m *Manager) RollDice(playerID string) (int, error) {
	if !m.IsPlayerTurn(playerID) {
		return 0, errors.New("сейчас ход другого игрока")
	}

	if m.Turn.DiceRolled {
		return 0, errors.New("кубик уже был брошен")
	}

	result := m.random.Intn(6) + 1

	m.Turn.DiceResult = result
	m.Turn.DiceRolled = true

	return result, nil
}

// Использование основного действия
func (m *Manager) UseAction(playerID string) error {
	if !m.IsPlayerTurn(playerID) {
		return errors.New("сейчас ход другого игрока")
	}

	if !m.Turn.DiceRolled {
		return errors.New("сначала необходимо бросить кубик")
	}

	if m.Turn.ActionUsed {
		return errors.New("действие за этот ход уже использовано")
	}

	m.Turn.ActionUsed = true

	return nil
}

// Завершение хода
func (m *Manager) EndTurn(playerID string) error {
	if !m.IsPlayerTurn(playerID) {
		return errors.New("сейчас ход другого игрока")
	}

	if !m.Turn.DiceRolled {
		return errors.New("сначала необходимо бросить кубик")
	}

	if !m.Turn.ActionUsed {
		return errors.New("необходимо выполнить основное действие")
	}

	// Определяем следующего игрока
	if m.Turn.CurrentPlayerID == m.PlayerIDs[0] {
		m.Turn.CurrentPlayerID = m.PlayerIDs[1]
	} else {
		m.Turn.CurrentPlayerID = m.PlayerIDs[0]
	}

	// Подготавливаем следующий ход
	m.Turn.Number++
	m.Turn.DiceResult = 0
	m.Turn.DiceRolled = false
	m.Turn.ActionUsed = false

	return nil
}
