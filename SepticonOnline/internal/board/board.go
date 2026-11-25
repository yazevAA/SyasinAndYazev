package board

import (
	"errors"
	"fmt"
)

const (
	// Количество уровней станции.
	LevelCount = 9

	// Количество подземных уровней.
	UndergroundLevels = 8

	// Количество отсеков на каждом подземном уровне.
	CompartmentsPerLevel = 21
)

var (
	ErrCompartmentNotFound    = errors.New("отсек не найден")
	ErrInvalidCompartmentType = errors.New("недопустимый тип отсека")
	ErrInvalidBoard           = errors.New("некорректная структура станции")
)

// Типы отсеков.
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

// Проверка допустимости типа отсека.
func IsValidCompartmentType(compartmentType CompartmentType) bool {
	switch compartmentType {
	case CompartmentUnknown,
		CompartmentProduction,
		CompartmentCombat,
		CompartmentEnergy,
		CompartmentResource,
		CompartmentRocket,
		CompartmentCommand:
		return true
	default:
		return false
	}
}

// Отдельный отсек станции.
type Compartment struct {
	ID       string          `json:"id"`
	Level    int             `json:"level"`
	Position int             `json:"position"`
	Type     CompartmentType `json:"type"`
	Damaged  bool            `json:"damaged"`
}

// Проверка, может ли отсек использоваться.
func (c Compartment) IsOperational() bool {
	return !c.Damaged && c.Type != CompartmentUnknown
}

// Подземный уровень станции.
type Level struct {
	Number       int           `json:"number"`
	Compartments []Compartment `json:"compartments"`
}

// Поверхность астероида.
type Surface struct {
	Level   int  `json:"level"`
	Present bool `json:"present"`
}

// Игровое поле одной станции.
type Board struct {
	Levels  []Level `json:"levels"`
	Surface Surface `json:"surface"`
}

// Создание нового игрового поля.
func NewBoard() Board {
	board := Board{
		Levels: make([]Level, 0, UndergroundLevels),
		Surface: Surface{
			Level:   LevelCount,
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

			level.Compartments = append(level.Compartments, compartment)
		}

		board.Levels = append(board.Levels, level)
	}

	return board
}

// Настройка типа отсека по его идентификатору.
func (b *Board) SetCompartmentType(id string, compartmentType CompartmentType) error {
	if !IsValidCompartmentType(compartmentType) {
		return ErrInvalidCompartmentType
	}

	compartment, err := b.FindCompartment(id)
	if err != nil {
		return err
	}

	compartment.Type = compartmentType
	return nil
}

// Проверка доступности отсека для игровых действий.
func (b *Board) IsCompartmentOperational(id string) (bool, error) {
	compartment, err := b.FindCompartment(id)
	if err != nil {
		return false, err
	}

	return compartment.IsOperational(), nil
}

// Проверка целостности структуры станции.
func (b *Board) Validate() error {
	if len(b.Levels) != UndergroundLevels {
		return fmt.Errorf(
			"%w: ожидалось %d уровней, получено %d",
			ErrInvalidBoard,
			UndergroundLevels,
			len(b.Levels),
		)
	}

	if b.Surface.Level != LevelCount {
		return fmt.Errorf(
			"%w: поверхность должна находиться на уровне %d",
			ErrInvalidBoard,
			LevelCount,
		)
	}

	for i, level := range b.Levels {
		expectedLevel := i + 1

		if level.Number != expectedLevel {
			return fmt.Errorf(
				"%w: ожидался уровень %d, получен %d",
				ErrInvalidBoard,
				expectedLevel,
				level.Number,
			)
		}

		if len(level.Compartments) != CompartmentsPerLevel {
			return fmt.Errorf(
				"%w: уровень %d должен содержать %d отсеков",
				ErrInvalidBoard,
				level.Number,
				CompartmentsPerLevel,
			)
		}

		for j, compartment := range level.Compartments {
			expectedPosition := j + 1
			expectedID := fmt.Sprintf("L%d-C%d", level.Number, expectedPosition)

			if compartment.ID != expectedID {
				return fmt.Errorf(
					"%w: неверный ID отсека на уровне %d, позиции %d",
					ErrInvalidBoard,
					level.Number,
					expectedPosition,
				)
			}

			if compartment.Level != level.Number {
				return fmt.Errorf(
					"%w: неверный уровень отсека %s",
					ErrInvalidBoard,
					compartment.ID,
				)
			}

			if compartment.Position != expectedPosition {
				return fmt.Errorf(
					"%w: неверная позиция отсека %s",
					ErrInvalidBoard,
					compartment.ID,
				)
			}

			if !IsValidCompartmentType(compartment.Type) {
				return fmt.Errorf(
					"%w: отсек %s",
					ErrInvalidCompartmentType,
					compartment.ID,
				)
			}
		}
	}

	return nil
}
