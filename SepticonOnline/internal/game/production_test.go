package game

import (
	"testing"

	"septikon-online/internal/board"
	"septikon-online/internal/resources"
)

// Проверка полного цикла производства ресурсов.
func TestProductionLifecycle(t *testing.T) {
	g := createProductionTestGame(t)

	player, err := g.FindPlayer("player1")
	if err != nil {
		t.Fatal(err)
	}

	if err := player.Station.SetCompartmentType(
		"L1-C1",
		board.CompartmentProduction,
	); err != nil {
		t.Fatal(err)
	}

	initialMetal := player.Inventory.Metal

	// Игрок бросает кубик.
	if _, err := g.RollDice("player1"); err != nil {
		t.Fatal(err)
	}

	// Запускаем производство 5 единиц металла за 2 хода.
	err = g.StartProduction(
		"player1",
		"production-1",
		"L1-C1",
		resources.Metal,
		5,
		2,
	)
	if err != nil {
		t.Fatalf("не удалось запустить производство: %v", err)
	}

	if len(player.ProductionTasks) != 1 {
		t.Fatalf("ожидалась одна задача производства")
	}

	task := &player.ProductionTasks[0]

	if task.RemainingTurns != 2 {
		t.Errorf("ожидалось 2 оставшихся хода, получено: %d", task.RemainingTurns)
	}

	if task.Completed {
		t.Error("производство не должно быть завершено сразу после запуска")
	}

	// Первый ход игрока завершается.
	if err := g.EndTurn("player1"); err != nil {
		t.Fatal(err)
	}

	if task.RemainingTurns != 1 {
		t.Errorf("ожидался 1 оставшийся ход, получено: %d", task.RemainingTurns)
	}

	// Второй игрок выполняет действие и завершает ход.
	if _, err := g.RollDice("player2"); err != nil {
		t.Fatal(err)
	}

	if err := g.TurnManager.UseAction("player2"); err != nil {
		t.Fatal(err)
	}

	if err := g.EndTurn("player2"); err != nil {
		t.Fatal(err)
	}

	if !task.Completed {
		t.Error("производство должно быть завершено после второго хода")
	}

	// Получаем произведённый металл.
	if _, err := g.RollDice("player1"); err != nil {
		t.Fatal(err)
	}

	if err := g.CollectProduction("player1", "production-1"); err != nil {
		t.Fatalf("не удалось получить ресурсы: %v", err)
	}

	if player.Inventory.Metal != initialMetal+5 {
		t.Errorf(
			"ожидалось %d металла, получено %d",
			initialMetal+5,
			player.Inventory.Metal,
		)
	}

	if !task.Collected {
		t.Error("задача должна быть отмечена как полученная")
	}

	if !task.Completed {
		t.Error("завершённое производство должно оставаться завершённым")
	}

	if task.Amount != 5 {
		t.Errorf("ожидалось сохранение количества 5, получено %d", task.Amount)
	}

	// Повторно получить те же ресурсы нельзя.
	metalBefore := player.Inventory.Metal

	if err := g.CollectProduction("player1", "production-1"); err == nil {
		t.Fatal("ожидалась ошибка при повторном получении ресурсов")
	}

	if player.Inventory.Metal != metalBefore {
		t.Error("при повторном получении количество ресурсов не должно изменяться")
	}
}

// Нельзя получить ресурсы до завершения производства.
func TestCollectProductionBeforeCompletion(t *testing.T) {
	g := createProductionTestGame(t)

	player, err := g.FindPlayer("player1")
	if err != nil {
		t.Fatal(err)
	}

	if err := player.Station.SetCompartmentType(
		"L1-C1",
		board.CompartmentProduction,
	); err != nil {
		t.Fatal(err)
	}

	if _, err := g.RollDice("player1"); err != nil {
		t.Fatal(err)
	}

	if err := g.StartProduction(
		"player1",
		"production-1",
		"L1-C1",
		resources.Metal,
		5,
		3,
	); err != nil {
		t.Fatal(err)
	}

	if err := g.EndTurn("player1"); err != nil {
		t.Fatal(err)
	}

	if _, err := g.RollDice("player2"); err != nil {
		t.Fatal(err)
	}

	if err := g.TurnManager.UseAction("player2"); err != nil {
		t.Fatal(err)
	}

	if err := g.EndTurn("player2"); err != nil {
		t.Fatal(err)
	}

	if _, err := g.RollDice("player1"); err != nil {
		t.Fatal(err)
	}

	err = g.CollectProduction("player1", "production-1")

	if err == nil {
		t.Fatal("ожидалась ошибка при попытке получить незавершённые ресурсы")
	}
}

// Нельзя запустить производство с некорректными параметрами.
func TestStartProductionInvalidParameters(t *testing.T) {
	tests := []struct {
		name          string
		resource      resources.ResourceType
		amount        int
		requiredTurns int
	}{
		{
			name:          "нулевое количество",
			resource:      resources.Metal,
			amount:        0,
			requiredTurns: 2,
		},
		{
			name:          "отрицательное количество",
			resource:      resources.Metal,
			amount:        -5,
			requiredTurns: 2,
		},
		{
			name:          "нулевая длительность",
			resource:      resources.Metal,
			amount:        5,
			requiredTurns: 0,
		},
		{
			name:          "неизвестный ресурс",
			resource:      "unknown",
			amount:        5,
			requiredTurns: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := createProductionTestGame(t)

			if _, err := g.RollDice("player1"); err != nil {
				t.Fatal(err)
			}

			err := g.StartProduction(
				"player1",
				"production-test",
				"L1-C1",
				tt.resource,
				tt.amount,
				tt.requiredTurns,
			)

			if err == nil {
				t.Fatal("ожидалась ошибка при некорректных параметрах")
			}
		})
	}
}
