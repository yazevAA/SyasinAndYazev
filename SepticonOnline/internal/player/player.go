package player

import "septikon-online/internal/board"

// Информация об игроке
type Player struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Ready   bool        `json:"ready"`
	Station board.Board `json:"station"`
}
