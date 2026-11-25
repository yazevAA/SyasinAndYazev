package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"septikon-online/internal/room"
)

// Создаёт тестовый HTTP-обработчик с новым менеджером комнат.
func newTestHandler() *Handler {
	manager := room.NewManager()
	return NewHandler(manager)
}

// Выполняет HTTP-запрос к обработчику.
func performRequest(
	handler http.Handler,
	method string,
	path string,
	body any,
) *httptest.ResponseRecorder {
	var requestBody bytes.Buffer

	if body != nil {
		_ = json.NewEncoder(&requestBody).Encode(body)
	}

	request := httptest.NewRequest(method, path, &requestBody)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder
}

// Проверяет создание комнаты.
func TestCreateRoom(t *testing.T) {
	handler := newTestHandler()

	body := CreateRoomRequest{
		PlayerID:   "player1",
		PlayerName: "Игрок 1",
	}

	response := performRequest(
		handler,
		http.MethodPost,
		"/api/rooms",
		body,
	)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusCreated,
			response.Code,
		)
	}

	if !bytes.Contains(response.Body.Bytes(), []byte("room-1")) {
		t.Fatal("в ответе отсутствует ID созданной комнаты")
	}
}

// Проверяет получение списка комнат.
func TestListRooms(t *testing.T) {
	handler := newTestHandler()

	performRequest(
		handler,
		http.MethodPost,
		"/api/rooms",
		CreateRoomRequest{
			PlayerID:   "player1",
			PlayerName: "Игрок 1",
		},
	)

	response := performRequest(
		handler,
		http.MethodGet,
		"/api/rooms",
		nil,
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusOK,
			response.Code,
		)
	}

	if !json.Valid(response.Body.Bytes()) {
		t.Fatal("сервер вернул некорректный JSON")
	}

	if !bytes.Contains(response.Body.Bytes(), []byte("room-1")) {
		t.Fatal("список комнат не содержит созданную комнату")
	}
}

// Проверяет получение информации о комнате.
func TestGetRoom(t *testing.T) {
	handler := newTestHandler()

	performRequest(
		handler,
		http.MethodPost,
		"/api/rooms",
		CreateRoomRequest{
			PlayerID:   "player1",
			PlayerName: "Игрок 1",
		},
	)

	response := performRequest(
		handler,
		http.MethodGet,
		"/api/rooms/room-1",
		nil,
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusOK,
			response.Code,
		)
	}

	if !json.Valid(response.Body.Bytes()) {
		t.Fatal("сервер вернул некорректный JSON")
	}
}

// Проверяет подключение второго игрока.
func TestJoinRoom(t *testing.T) {
	handler := newTestHandler()

	performRequest(
		handler,
		http.MethodPost,
		"/api/rooms",
		CreateRoomRequest{
			PlayerID:   "player1",
			PlayerName: "Игрок 1",
		},
	)

	response := performRequest(
		handler,
		http.MethodPost,
		"/api/rooms/room-1/join",
		JoinRoomRequest{
			PlayerID:   "player2",
			PlayerName: "Игрок 2",
		},
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"ожидался статус %d, получен %d. Ответ: %s",
			http.StatusOK,
			response.Code,
			response.Body.String(),
		)
	}

	if !bytes.Contains(response.Body.Bytes(), []byte("player2")) {
		t.Fatal("в ответе отсутствует второй игрок")
	}
}

// Проверяет бросок кубика.
func TestRollDice(t *testing.T) {
	handler := newTestHandler()

	performRequest(
		handler,
		http.MethodPost,
		"/api/rooms",
		CreateRoomRequest{
			PlayerID:   "player1",
			PlayerName: "Игрок 1",
		},
	)

	performRequest(
		handler,
		http.MethodPost,
		"/api/rooms/room-1/join",
		JoinRoomRequest{
			PlayerID:   "player2",
			PlayerName: "Игрок 2",
		},
	)

	response := performRequest(
		handler,
		http.MethodPost,
		"/api/rooms/room-1/game/roll",
		map[string]string{
			"playerId": "player1",
		},
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"ожидался статус %d, получен %d. Ответ: %s",
			http.StatusOK,
			response.Code,
			response.Body.String(),
		)
	}

	var dice DiceResponse

	if err := json.Unmarshal(response.Body.Bytes(), &dice); err != nil {
		t.Fatalf("не удалось прочитать результат кубика: %v", err)
	}

	if dice.Result < 1 || dice.Result > 6 {
		t.Fatalf("недопустимый результат кубика: %d", dice.Result)
	}
}

// Проверяет обращение к несуществующей комнате.
func TestGetNonexistentRoom(t *testing.T) {
	handler := newTestHandler()

	response := performRequest(
		handler,
		http.MethodGet,
		"/api/rooms/nonexistent",
		nil,
	)

	if response.Code != http.StatusNotFound {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusNotFound,
			response.Code,
		)
	}
}

// Проверяет обработку некорректного JSON.
func TestInvalidJSON(t *testing.T) {
	handler := newTestHandler()

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/rooms",
		bytes.NewBufferString("{invalid json}"),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusBadRequest,
			response.Code,
		)
	}
}

// Проверяет запрет неподдерживаемого HTTP-метода.
func TestMethodNotAllowed(t *testing.T) {
	handler := newTestHandler()

	response := performRequest(
		handler,
		http.MethodDelete,
		"/api/rooms",
		nil,
	)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusMethodNotAllowed,
			response.Code,
		)
	}
}

func TestMoveRobot(t *testing.T) {
	manager := room.NewManager()
	handler := NewHandler(manager)

	// Создаём комнату.
	_, err := manager.CreateRoom("player1", "Игрок 1")
	if err != nil {
		t.Fatalf("не удалось создать комнату: %v", err)
	}

	// Подключаем второго игрока.
	_, err = manager.JoinRoom("room-1", "player2", "Игрок 2")
	if err != nil {
		t.Fatalf("не удалось подключить второго игрока: %v", err)
	}

	// Получаем комнату и добавляем робота первому игроку.
	currentRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить комнату: %v", err)
	}

	err = currentRoom.Game.AddRobot("player1", "robot-1", 1, 1)
	if err != nil {
		t.Fatalf("не удалось добавить робота: %v", err)
	}

	// Бросаем кубик через HTTP API.
	rollRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/rooms/room-1/game/roll",
		bytes.NewBufferString(`{"playerId":"player1"}`),
	)

	rollResponse := httptest.NewRecorder()
	handler.ServeHTTP(rollResponse, rollRequest)

	if rollResponse.Code != http.StatusOK {
		t.Fatalf(
			"ожидался статус 200 при броске кубика, получен %d: %s",
			rollResponse.Code,
			rollResponse.Body.String(),
		)
	}

	var dice DiceResponse

	if err := json.NewDecoder(rollResponse.Body).Decode(&dice); err != nil {
		t.Fatalf("не удалось прочитать результат кубика: %v", err)
	}

	if dice.Result < 1 || dice.Result > 6 {
		t.Fatalf("некорректный результат кубика: %d", dice.Result)
	}

	// Отправляем робота по маршруту.
	// Маршрут содержит 7 позиций, поэтому подходит для любого результата кубика.
	moveRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/rooms/room-1/game/move-robot",
		bytes.NewBufferString(`{
			"playerId": "player1",
			"robotId": "robot-1",
			"path": [
				{"level": 1, "position": 1},
				{"level": 1, "position": 2},
				{"level": 1, "position": 3},
				{"level": 1, "position": 4},
				{"level": 1, "position": 5},
				{"level": 1, "position": 6},
				{"level": 1, "position": 7}
			]
		}`),
	)

	moveResponse := httptest.NewRecorder()
	handler.ServeHTTP(moveResponse, moveRequest)

	if moveResponse.Code != http.StatusOK {
		t.Fatalf(
			"ожидался статус 200 при перемещении робота, получен %d: %s",
			moveResponse.Code,
			moveResponse.Body.String(),
		)
	}

	// Проверяем фактическую позицию робота.
	updatedRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить обновлённую комнату: %v", err)
	}

	robot := updatedRoom.Game.Players[0].Robots[0]
	expectedPosition := 1 + dice.Result

	if robot.Position.Position != expectedPosition {
		t.Fatalf(
			"робот переместился неправильно: ожидалась позиция %d, получена %d",
			expectedPosition,
			robot.Position.Position,
		)
	}

	if robot.Position.Level != 1 {
		t.Fatalf(
			"ожидался уровень 1, получен уровень %d",
			robot.Position.Level,
		)
	}

	t.Logf(
		"робот успешно перемещён на %d клеток, новая позиция: уровень %d, клетка %d",
		dice.Result,
		robot.Position.Level,
		robot.Position.Position,
	)
}

func newRepairTestEnvironment(t *testing.T, metal int) (*room.Manager, *Handler) {
	t.Helper()

	manager := room.NewManager()
	handler := NewHandler(manager)

	_, err := manager.CreateRoom("player1", "Игрок 1")
	if err != nil {
		t.Fatalf("не удалось создать комнату: %v", err)
	}

	_, err = manager.JoinRoom("room-1", "player2", "Игрок 2")
	if err != nil {
		t.Fatalf("не удалось подключить второго игрока: %v", err)
	}

	currentRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить комнату: %v", err)
	}

	// Настраиваем количество металла у первого игрока.
	currentRoom.Game.Players[0].Inventory.Metal = metal

	return manager, handler
}

func prepareRepairTest(t *testing.T, handler *Handler) {
	t.Helper()

	// Бросаем кубик перед началом ремонта.
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/rooms/room-1/game/roll",
		bytes.NewBufferString(`{"playerId":"player1"}`),
	)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"не удалось бросить кубик: статус %d, ответ: %s",
			response.Code,
			response.Body.String(),
		)
	}
}

func TestStartRepair(t *testing.T) {
	manager, handler := newRepairTestEnvironment(t, 10)

	currentRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить комнату: %v", err)
	}

	// Повреждаем отсек первого игрока.
	err = currentRoom.Game.Players[0].Station.DamageCompartment("L1-C1")
	if err != nil {
		t.Fatalf("не удалось повредить отсек: %v", err)
	}

	prepareRepairTest(t, handler)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/rooms/room-1/game/repair",
		bytes.NewBufferString(`{
			"playerId": "player1",
			"taskId": "repair-1",
			"compartmentId": "L1-C1"
		}`),
	)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"ожидался статус 200 при запуске ремонта, получен %d: %s",
			response.Code,
			response.Body.String(),
		)
	}

	updatedRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить обновлённую комнату: %v", err)
	}

	player := updatedRoom.Game.Players[0]

	// Проверяем списание металла.
	if player.Inventory.Metal != 8 {
		t.Fatalf(
			"ожидалось 8 единиц металла, получено %d",
			player.Inventory.Metal,
		)
	}

	// Проверяем созданную задачу.
	task, err := updatedRoom.Game.GetRepairStatus("repair-1")
	if err != nil {
		t.Fatalf("задача ремонта не создана: %v", err)
	}

	if task.OwnerID != "player1" {
		t.Errorf("неверный владелец ремонта: %s", task.OwnerID)
	}

	if task.CompartmentID != "L1-C1" {
		t.Errorf("неверный ID отсека: %s", task.CompartmentID)
	}

	if task.RemainingTurns != 2 {
		t.Errorf(
			"ожидалась длительность ремонта 2 хода, получено %d",
			task.RemainingTurns,
		)
	}

	if task.Completed {
		t.Error("новая задача ремонта не должна быть завершена")
	}

	t.Log("ремонт успешно запущен, металл списан, задача создана")
}

func TestRepairInsufficientMetal(t *testing.T) {
	manager, handler := newRepairTestEnvironment(t, 1)

	currentRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить комнату: %v", err)
	}

	err = currentRoom.Game.Players[0].Station.DamageCompartment("L1-C1")
	if err != nil {
		t.Fatalf("не удалось повредить отсек: %v", err)
	}

	prepareRepairTest(t, handler)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/rooms/room-1/game/repair",
		bytes.NewBufferString(`{
			"playerId": "player1",
			"taskId": "repair-2",
			"compartmentId": "L1-C1"
		}`),
	)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code == http.StatusOK {
		t.Fatal("ремонт не должен запускаться при недостатке металла")
	}

	updatedRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить комнату: %v", err)
	}

	player := updatedRoom.Game.Players[0]

	if player.Inventory.Metal != 1 {
		t.Errorf(
			"металл не должен списываться при ошибке, осталось: %d",
			player.Inventory.Metal,
		)
	}

	if _, err := updatedRoom.Game.GetRepairStatus("repair-2"); err == nil {
		t.Error("задача ремонта не должна была создаться")
	}

	t.Log("недостаток металла корректно обработан")
}

func TestRepairUndamagedCompartment(t *testing.T) {
	manager, handler := newRepairTestEnvironment(t, 10)

	// Отсек L1-C1 не повреждаем.
	prepareRepairTest(t, handler)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/rooms/room-1/game/repair",
		bytes.NewBufferString(`{
			"playerId": "player1",
			"taskId": "repair-3",
			"compartmentId": "L1-C1"
		}`),
	)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code == http.StatusOK {
		t.Fatal("нельзя ремонтировать неповреждённый отсек")
	}

	updatedRoom, err := manager.GetRoom("room-1")
	if err != nil {
		t.Fatalf("не удалось получить комнату: %v", err)
	}

	player := updatedRoom.Game.Players[0]

	if player.Inventory.Metal != 10 {
		t.Errorf(
			"металл не должен списываться, осталось: %d",
			player.Inventory.Metal,
		)
	}

	if _, err := updatedRoom.Game.GetRepairStatus("repair-3"); err == nil {
		t.Error("задача ремонта не должна была создаться")
	}

	t.Log("попытка ремонта исправного отсека корректно отклонена")
}
