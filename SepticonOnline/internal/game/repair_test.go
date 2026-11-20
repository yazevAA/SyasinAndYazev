package game

import "testing"

func TestRepair(t *testing.T) {
	g, err := NewGame(
		"game1",
		"player1",
		"Игрок 1",
		"player2",
		"Игрок 2",
	)

	if err != nil {
		t.Fatal(err)
	}

	// Повреждаем отсек
	err = g.Players[0].Station.DamageCompartment("L1-C1")

	if err != nil {
		t.Fatal(err)
	}

	// Бросаем кубик
	_, err = g.RollDice("player1")

	if err != nil {
		t.Fatal(err)
	}

	// Начинаем ремонт
	err = g.StartRepair(
		"player1",
		"repair1",
		"L1-C1",
	)

	if err != nil {
		t.Fatal(err)
	}

	// Первый игровой ход
	g.AdvanceRepairs()

	if g.IsRepairCompleted("repair1") {
		t.Error("ремонт завершился слишком рано")
	}

	// Второй игровой ход
	g.AdvanceRepairs()

	if !g.IsRepairCompleted("repair1") {
		t.Error("ремонт должен быть завершён")
	}

	compartment, err := g.Players[0].Station.FindCompartment("L1-C1")

	if err != nil {
		t.Fatal(err)
	}

	if compartment.Damaged {
		t.Error("отсек должен быть восстановлен")
	}
}
