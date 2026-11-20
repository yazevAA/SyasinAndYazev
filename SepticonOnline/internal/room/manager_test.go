package room

import "testing"

func TestRoomCreation(t *testing.T) {
	manager := NewManager()

	createdRoom, err := manager.CreateRoom("player1", "Игрок 1")
	if err != nil {
		t.Fatal(err)
	}

	if createdRoom.Status != RoomWaiting {
		t.Error("комната должна ожидать второго игрока")
	}

	if len(createdRoom.Players) != 1 {
		t.Error("в комнате должен находиться один игрок")
	}

	if createdRoom.Game != nil {
		t.Error("игра не должна быть создана до подключения второго игрока")
	}
}

func TestJoinRoom(t *testing.T) {
	manager := NewManager()

	createdRoom, err := manager.CreateRoom("player1", "Игрок 1")
	if err != nil {
		t.Fatal(err)
	}

	joinedRoom, err := manager.JoinRoom(
		createdRoom.ID,
		"player2",
		"Игрок 2",
	)
	if err != nil {
		t.Fatal(err)
	}

	if joinedRoom.Status != RoomPlaying {
		t.Error("комната должна перейти в состояние игры")
	}

	if len(joinedRoom.Players) != 2 {
		t.Error("в комнате должно быть два игрока")
	}

	if joinedRoom.Game == nil {
		t.Fatal("игровая партия должна быть создана")
	}

	if len(joinedRoom.Game.Players) != 2 {
		t.Error("в игровой партии должно быть два игрока")
	}
}

func TestJoinFullRoom(t *testing.T) {
	manager := NewManager()

	createdRoom, err := manager.CreateRoom("player1", "Игрок 1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.JoinRoom(createdRoom.ID, "player2", "Игрок 2")
	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.JoinRoom(createdRoom.ID, "player3", "Игрок 3")
	if err == nil {
		t.Error("нельзя подключить третьего игрока")
	}
}
