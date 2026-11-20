package turns

import "testing"

// Проверка создания системы ходов
func TestNewManager(t *testing.T) {
	manager, err := NewManager([]string{"player1", "player2"})

	if err != nil {
		t.Fatal(err)
	}

	if manager.Turn.CurrentPlayerID != "player1" {
		t.Error("неправильно определён первый игрок")
	}
}

// Проверка броска кубика
func TestRollDice(t *testing.T) {
	manager, _ := NewManager([]string{"player1", "player2"})

	result, err := manager.RollDice("player1")

	if err != nil {
		t.Fatal(err)
	}

	if result < 1 || result > 6 {
		t.Errorf("некорректный результат кубика: %d", result)
	}
}

// Проверка запрета повторного броска
func TestSecondRoll(t *testing.T) {
	manager, _ := NewManager([]string{"player1", "player2"})

	_, err := manager.RollDice("player1")

	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.RollDice("player1")

	if err == nil {
		t.Error("повторный бросок должен быть запрещён")
	}
}

// Проверка запрета второго действия
func TestSecondAction(t *testing.T) {
	manager, _ := NewManager([]string{"player1", "player2"})

	_, _ = manager.RollDice("player1")

	err := manager.UseAction("player1")

	if err != nil {
		t.Fatal(err)
	}

	err = manager.UseAction("player1")

	if err == nil {
		t.Error("второе действие должно быть запрещено")
	}
}

// Проверка переключения игрока
func TestEndTurn(t *testing.T) {
	manager, _ := NewManager([]string{"player1", "player2"})

	_, _ = manager.RollDice("player1")
	_ = manager.UseAction("player1")

	err := manager.EndTurn("player1")

	if err != nil {
		t.Fatal(err)
	}

	if manager.Turn.CurrentPlayerID != "player2" {
		t.Error("ход не перешёл ко второму игроку")
	}

	if manager.Turn.Number != 2 {
		t.Error("номер хода не увеличился")
	}
}
