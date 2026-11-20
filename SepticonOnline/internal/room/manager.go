package room

import (
	"errors"
	"fmt"
	"septikon-online/internal/game"
	"septikon-online/internal/robot"
	"sync"
)

var ErrRoomNotFound = errors.New("комната не найдена")

type Manager struct {
	rooms  map[string]*Room
	mutex  sync.Mutex
	nextID uint64
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

// Создание комнаты первым игроком
func (m *Manager) CreateRoom(playerID, playerName string) (Room, error) {
	if playerID == "" || playerName == "" {
		return Room{}, errors.New("ID и имя игрока не могут быть пустыми")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.nextID++
	roomID := fmt.Sprintf("room-%d", m.nextID)

	newRoom := &Room{
		ID:     roomID,
		Status: RoomWaiting,
		Players: []RoomPlayer{
			{
				ID:   playerID,
				Name: playerName,
			},
		},
	}

	m.rooms[roomID] = newRoom

	return *newRoom, nil
}

// Подключение второго игрока
func (m *Manager) JoinRoom(roomID, playerID, playerName string) (Room, error) {
	if playerID == "" || playerName == "" {
		return Room{}, errors.New("ID и имя игрока не могут быть пустыми")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	currentRoom, exists := m.rooms[roomID]
	if !exists {
		return Room{}, ErrRoomNotFound
	}

	if currentRoom.Status != RoomWaiting {
		return Room{}, errors.New("комната уже заполнена")
	}

	if len(currentRoom.Players) >= MaxPlayers {
		return Room{}, errors.New("комната заполнена")
	}

	if currentRoom.Players[0].ID == playerID {
		return Room{}, errors.New("игрок уже находится в комнате")
	}

	firstPlayer := currentRoom.Players[0]

	// Создаём общую игровую партию
	newGame, err := game.NewGame(
		currentRoom.ID+"-game",
		firstPlayer.ID,
		firstPlayer.Name,
		playerID,
		playerName,
	)

	if err != nil {
		return Room{}, err
	}

	currentRoom.Players = append(currentRoom.Players, RoomPlayer{
		ID:   playerID,
		Name: playerName,
	})

	currentRoom.Game = newGame
	currentRoom.Status = RoomPlaying

	return *currentRoom, nil
}

// Получение информации о комнате
func (m *Manager) GetRoom(roomID string) (Room, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	currentRoom, exists := m.rooms[roomID]
	if !exists {
		return Room{}, ErrRoomNotFound
	}

	return *currentRoom, nil
}

func (m *Manager) ListRooms() []Room {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	result := make([]Room, 0, len(m.rooms))

	for _, currentRoom := range m.rooms {
		// Показываем только комнаты, ожидающие второго игрока.
		if currentRoom.Status == RoomWaiting {
			result = append(result, *currentRoom)
		}
	}

	return result
}

// Бросок кубика игроком
func (m *Manager) RollDice(roomID, playerID string) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	currentRoom, exists := m.rooms[roomID]
	if !exists {
		return 0, ErrRoomNotFound
	}

	if currentRoom.Status != RoomPlaying || currentRoom.Game == nil {
		return 0, errors.New("партия еще не началась")
	}

	return currentRoom.Game.RollDice(playerID)
}

// Завершение хода игроком
func (m *Manager) EndTurn(roomID, playerID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	currentRoom, exists := m.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	if currentRoom.Status != RoomPlaying || currentRoom.Game == nil {
		return errors.New("партия еще не началась")
	}

	return currentRoom.Game.EndTurn(playerID)
}

// Перемещение робота
func (m *Manager) MoveRobot(
	roomID string,
	playerID string,
	robotID string,
	path []robot.Position,
) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	currentRoom, exists := m.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	if currentRoom.Status != RoomPlaying || currentRoom.Game == nil {
		return errors.New("партия еще не началась")
	}

	return currentRoom.Game.MoveRobot(playerID, robotID, path)
}

// Начало ремонта отсека
func (m *Manager) StartRepair(
	roomID string,
	playerID string,
	taskID string,
	compartmentID string,
) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	currentRoom, exists := m.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	if currentRoom.Status != RoomPlaying || currentRoom.Game == nil {
		return errors.New("партия еще не началась")
	}

	return currentRoom.Game.StartRepair(playerID, taskID, compartmentID)
}
