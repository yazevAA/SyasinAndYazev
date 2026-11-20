package game

import (
	"errors"

	"septikon-online/internal/board"
	"septikon-online/internal/resources"
)

// Стоимость и длительность ремонта
const (
	RepairMetalCost = 2
	RepairDuration  = 2
)

// Задача ремонта
type RepairTask struct {
	ID             string `json:"id"`
	OwnerID        string `json:"ownerId"`
	CompartmentID  string `json:"compartmentId"`
	RemainingTurns int    `json:"remainingTurns"`
	Completed      bool   `json:"completed"`
}

// Начало ремонта
func (g *Game) StartRepair(
	playerID string,
	taskID string,
	compartmentID string,
) error {

	if err := g.CheckTurn(playerID); err != nil {
		return err
	}

	if !g.TurnManager.Turn.DiceRolled {
		return errors.New("сначала необходимо бросить кубик")
	}

	if g.TurnManager.Turn.ActionUsed {
		return errors.New("действие уже использовано")
	}

	if taskID == "" {
		return errors.New("ID ремонта не может быть пустым")
	}

	p, err := g.FindPlayer(playerID)
	if err != nil {
		return err
	}

	for _, task := range g.Repairs {
		if task.ID == taskID {
			return errors.New("ремонт с таким ID уже существует")
		}

		if task.OwnerID == playerID &&
			task.CompartmentID == compartmentID &&
			!task.Completed {
			return errors.New("этот отсек уже ремонтируется")
		}
	}

	compartment, err := p.Station.FindCompartment(compartmentID)
	if err != nil {
		return err
	}

	if !compartment.Damaged {
		return errors.New("отсек не повреждён")
	}

	if err := p.Inventory.Spend(
		resources.Metal,
		RepairMetalCost,
	); err != nil {
		return err
	}

	task := RepairTask{
		ID:             taskID,
		OwnerID:        playerID,
		CompartmentID:  compartmentID,
		RemainingTurns: RepairDuration,
		Completed:      false,
	}

	g.Repairs = append(g.Repairs, task)

	return g.TurnManager.UseAction(playerID)
}

// Продвижение ремонта
func (g *Game) AdvanceRepairs() {
	for i := range g.Repairs {
		task := &g.Repairs[i]

		if task.Completed {
			continue
		}

		task.RemainingTurns--

		if task.RemainingTurns <= 0 {
			task.RemainingTurns = 0
			task.Completed = true

			p, err := g.FindPlayer(task.OwnerID)
			if err != nil {
				continue
			}

			compartment, err := p.Station.FindCompartment(
				task.CompartmentID,
			)

			if err == nil {
				compartment.Damaged = false
			}
		}
	}
}

// Получение статуса ремонта
func (g *Game) GetRepairStatus(taskID string) (*RepairTask, error) {
	for i := range g.Repairs {
		if g.Repairs[i].ID == taskID {
			return &g.Repairs[i], nil
		}
	}

	return nil, errors.New("задача ремонта не найдена")
}

// Проверка завершения ремонта
func (g *Game) IsRepairCompleted(taskID string) bool {
	task, err := g.GetRepairStatus(taskID)

	if err != nil {
		return false
	}

	return task.Completed
}

// Обеспечение совместимости с моделью станции
var _ *board.Compartment
