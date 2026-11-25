package game

import (
	"testing"

	"septikon-online/internal/board"
	"septikon-online/internal/resources"
)

func createProductionTestGame(t *testing.T) *Game {
	t.Helper()

	g, err := NewGame(
		"test-game",
		"player1",
		"Player 1",
		"player2",
		"Player 2",
	)
	if err != nil {
		t.Fatal(err)
	}

	return g
}

func prepareProductionTurn(t *testing.T, g *Game) {
	t.Helper()

	_, err := g.RollDice("player1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartProductionInCorrectCompartment(t *testing.T) {
	g := createProductionTestGame(t)

	p, err := g.FindPlayer("player1")
	if err != nil {
		t.Fatal(err)
	}

	err = p.Station.SetCompartmentType(
		"L1-C1",
		board.CompartmentProduction,
	)
	if err != nil {
		t.Fatal(err)
	}

	prepareProductionTurn(t, g)

	err = g.StartProduction(
		"player1",
		"task-1",
		"L1-C1",
		resources.Metal,
		5,
		2,
	)
	if err != nil {
		t.Fatalf("производство должно запуститься: %v", err)
	}

	if len(p.ProductionTasks) != 1 {
		t.Fatal("задача производства не была создана")
	}

	if p.ProductionTasks[0].CompartmentID != "L1-C1" {
		t.Error("задача должна быть привязана к L1-C1")
	}
}

func TestStartProductionInWrongCompartmentType(t *testing.T) {
	g := createProductionTestGame(t)

	p, _ := g.FindPlayer("player1")

	err := p.Station.SetCompartmentType(
		"L1-C1",
		board.CompartmentEnergy,
	)
	if err != nil {
		t.Fatal(err)
	}

	prepareProductionTurn(t, g)

	err = g.StartProduction(
		"player1",
		"task-1",
		"L1-C1",
		resources.Metal,
		5,
		2,
	)

	if err == nil {
		t.Fatal("ожидалась ошибка: отсек не предназначен для производства")
	}

	if len(p.ProductionTasks) != 0 {
		t.Error("задача не должна создаваться")
	}

	if g.TurnManager.Turn.ActionUsed {
		t.Error("действие не должно считаться использованным")
	}
}

func TestStartProductionInDamagedCompartment(t *testing.T) {
	g := createProductionTestGame(t)

	p, _ := g.FindPlayer("player1")

	err := p.Station.SetCompartmentType(
		"L1-C1",
		board.CompartmentProduction,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Station.DamageCompartment("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	prepareProductionTurn(t, g)

	err = g.StartProduction(
		"player1",
		"task-1",
		"L1-C1",
		resources.Metal,
		5,
		2,
	)

	if err == nil {
		t.Fatal("ожидалась ошибка повреждённого отсека")
	}

	if len(p.ProductionTasks) != 0 {
		t.Error("задача не должна создаваться")
	}
}

func TestStartProductionInUnknownCompartment(t *testing.T) {
	g := createProductionTestGame(t)

	prepareProductionTurn(t, g)

	p, _ := g.FindPlayer("player1")

	err := g.StartProduction(
		"player1",
		"task-1",
		"L1-C1",
		resources.Metal,
		5,
		2,
	)

	if err == nil {
		t.Fatal("ожидалась ошибка для отсека unknown")
	}

	if len(p.ProductionTasks) != 0 {
		t.Error("задача не должна создаваться")
	}
}
