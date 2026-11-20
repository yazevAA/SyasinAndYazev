package room

import "septikon-online/internal/game"

const MaxPlayers = 2

type RoomStatus string

const (
	RoomWaiting RoomStatus = "waiting"
	RoomPlaying RoomStatus = "playing"
)

type RoomPlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Room struct {
	ID      string       `json:"id"`
	Status  RoomStatus   `json:"status"`
	Players []RoomPlayer `json:"players"`
	Game    *game.Game   `json:"game,omitempty"`
}
