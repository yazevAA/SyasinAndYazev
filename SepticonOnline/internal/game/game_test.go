package game

import "testing"

func TestNewGame(t *testing.T) {
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

	if len(g.Players) != 2 {
		t.Fatal("должно быть два игрока")
	}

	if g.Status != StatusPlaying {
		t.Error("неправильный статус партии")
	}

	if g.Players[0].ID != "player1" {
		t.Error("неправильно создан первый игрок")
	}

	if len(g.Players[0].Station.Levels) != 8 {
		t.Error("неправильное количество уровней станции")
	}
}

func TestGameTurn(t *testing.T) {
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

	_, err = g.RollDice("player1")

	if err != nil {
		t.Fatal(err)
	}

	err = g.TurnManager.UseAction("player1")

	if err != nil {
		t.Fatal(err)
	}

	err = g.EndTurn("player1")

	if err != nil {
		t.Fatal(err)
	}

	if g.TurnManager.Turn.CurrentPlayerID != "player2" {
		t.Error("ход должен перейти второму игроку")
	}
}
