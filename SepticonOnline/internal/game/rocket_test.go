package game

import (
	"testing"

	"septikon-online/internal/rocket"
)

// Создание тестовой игры с двумя игроками.
func createRocketTestGame(t *testing.T) *Game {
	t.Helper()

	g, err := NewGame(
		"test-game",
		"player-1",
		"Игрок 1",
		"player-2",
		"Игрок 2",
	)
	if err != nil {
		t.Fatalf("не удалось создать игру: %v", err)
	}

	return g
}

// Создание робота-стрелка на поверхности первой станции.
func addTestGunner(t *testing.T, g *Game) {
	t.Helper()

	if err := g.AddRobot("player-1", "gunner-1", 9, 1); err != nil {
		t.Fatalf("не удалось создать стрелка: %v", err)
	}
}

// Проверка успешного запуска ракеты.
func TestLaunchRocket(t *testing.T) {
	g := createRocketTestGame(t)

	if _, err := g.RollDice("player-1"); err != nil {
		t.Fatalf("не удалось бросить кубик: %v", err)
	}

	addTestGunner(t, g)

	err := g.LaunchRocket(
		"player-1",
		"rocket-1",
		"gunner-1",
		"player-2",
		10,
		2,
	)
	if err != nil {
		t.Fatalf("не удалось запустить ракету: %v", err)
	}

	if len(g.Rockets) != 1 {
		t.Fatalf("ожидалась одна ракета, получено: %d", len(g.Rockets))
	}

	r := g.Rockets[0]

	if r.Status != rocket.RocketFlying {
		t.Errorf("ожидался статус flying, получено: %s", r.Status)
	}

	if r.OwnerID != "player-1" {
		t.Errorf("неверный владелец ракеты: %s", r.OwnerID)
	}

	if r.TargetStation != "player-2" {
		t.Errorf("неверная целевая станция: %s", r.TargetStation)
	}

	if r.Position != 0 {
		t.Errorf("ракета не должна перемещаться до завершения хода")
	}
}

// Проверка продвижения ракеты после завершения хода.
func TestRocketAdvancesAfterEndTurn(t *testing.T) {
	g := createRocketTestGame(t)

	if _, err := g.RollDice("player-1"); err != nil {
		t.Fatal(err)
	}

	addTestGunner(t, g)

	err := g.LaunchRocket(
		"player-1",
		"rocket-1",
		"gunner-1",
		"player-2",
		10,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Первый игрок завершает ход.
	if err := g.EndTurn("player-1"); err != nil {
		t.Fatal(err)
	}

	if g.Rockets[0].Position != 2 {
		t.Errorf(
			"ожидалась позиция 2, получено: %d",
			g.Rockets[0].Position,
		)
	}

	// Второй игрок бросает кубик.
	if _, err := g.RollDice("player-2"); err != nil {
		t.Fatal(err)
	}

	// Выполняем основное действие для завершения хода.
	if err := g.TurnManager.UseAction("player-2"); err != nil {
		t.Fatal(err)
	}

	if err := g.EndTurn("player-2"); err != nil {
		t.Fatal(err)
	}

	// Ракета должна продвинуться ещё на 2 позиции.
	if g.Rockets[0].Position != 4 {
		t.Errorf(
			"ожидалась позиция 4, получено: %d",
			g.Rockets[0].Position,
		)
	}

	if g.Rockets[0].Status != rocket.RocketFlying {
		t.Errorf(
			"ожидался статус flying, получено: %s",
			g.Rockets[0].Status,
		)
	}
}

// Проверка запрета повторного использования ID ракеты.
func TestLaunchRocketDuplicateID(t *testing.T) {
	g := createRocketTestGame(t)

	if _, err := g.RollDice("player-1"); err != nil {
		t.Fatal(err)
	}

	addTestGunner(t, g)

	err := g.LaunchRocket(
		"player-1",
		"rocket-1",
		"gunner-1",
		"player-2",
		10,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := g.EndTurn("player-1"); err != nil {
		t.Fatal(err)
	}

	if _, err := g.RollDice("player-2"); err != nil {
		t.Fatal(err)
	}

	if err := g.TurnManager.UseAction("player-2"); err != nil {
		t.Fatal(err)
	}

	if err := g.EndTurn("player-2"); err != nil {
		t.Fatal(err)
	}

	if _, err := g.RollDice("player-1"); err != nil {
		t.Fatal(err)
	}

	err = g.LaunchRocket(
		"player-1",
		"rocket-1",
		"gunner-1",
		"player-2",
		10,
		2,
	)

	if err == nil {
		t.Fatal("ожидалась ошибка при повторном использовании ID ракеты")
	}
}
