package board

import "fmt"

// Количество уровней станции
const LevelCount = 9

// Количество подземных уровней
const UndergroundLevels = 8

// Количество отсеков на каждом подземном уровне
const CompartmentsPerLevel = 21

// Типы отсеков
type CompartmentType string

const (
	CompartmentUnknown    CompartmentType = "unknown"
	CompartmentProduction CompartmentType = "production"
	CompartmentCombat     CompartmentType = "combat"
	CompartmentEnergy     CompartmentType = "energy"
	CompartmentResource   CompartmentType = "resource"
	CompartmentRocket     CompartmentType = "rocket"
	CompartmentCommand    CompartmentType = "command"
)

// Отдельный отсек станции
type Compartment struct {
	ID       string          `json:"id"`
	Level    int             `json:"level"`
	Position int             `json:"position"`
	Type     CompartmentType `json:"type"`
	Damaged  bool            `json:"damaged"`
}

// Подземный уровень станции
type Level struct {
	Number       int           `json:"number"`
	Compartments []Compartment `json:"compartments"`
}

// Поверхность астероида
type Surface struct {
	Level   int  `json:"level"`
	Present bool `json:"present"`
}

// Игровое поле одной станции
type Board struct {
	Levels  []Level `json:"levels"`
	Surface Surface `json:"surface"`
}

// Создание нового игрового поля
func NewBoard() Board {
	board := Board{
		Levels: make([]Level, 0, UndergroundLevels),
		Surface: Surface{
			Level:   9,
			Present: true,
		},
	}

	for levelNumber := 1; levelNumber <= UndergroundLevels; levelNumber++ {
		level := Level{
			Number:       levelNumber,
			Compartments: make([]Compartment, 0, CompartmentsPerLevel),
		}

		for position := 1; position <= CompartmentsPerLevel; position++ {
			compartment := Compartment{
				ID:       fmt.Sprintf("L%d-C%d", levelNumber, position),
				Level:    levelNumber,
				Position: position,
				Type:     CompartmentUnknown,
				Damaged:  false,
			}

			level.Compartments = append(
				level.Compartments,
				compartment,
			)
		}

		board.Levels = append(board.Levels, level)
	}

	return board
}
