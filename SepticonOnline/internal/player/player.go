package player

import (
	"septikon-online/internal/board"
	"septikon-online/internal/resources"
	"septikon-online/internal/robot"
)

// Информация об игроке
type Player struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	Ready           bool                       `json:"ready"`
	Station         board.Board                `json:"station"`
	Inventory       resources.Inventory        `json:"inventory"`
	Robots          []robot.Robot              `json:"robots"`
	ProductionTasks []resources.ProductionTask `json:"productionTasks"`
}

// Создание игрока
func NewPlayer(id string, name string) Player {
	return Player{
		ID:      id,
		Name:    name,
		Ready:   false,
		Station: board.NewBoard(),
		Inventory: resources.Inventory{
			Metal:   10,
			Uranium: 5,
			Oxygen:  10,
			Biomass: 5,
			Energy:  10,
		},
		Robots:          make([]robot.Robot, 0),
		ProductionTasks: make([]resources.ProductionTask, 0),
	}
}
