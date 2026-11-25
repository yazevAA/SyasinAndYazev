package robot

import "errors"

// Положение робота на игровом поле
type Position struct {
	Level    int `json:"level"`
	Position int `json:"position"`
}

// Игровой робот
type Robot struct {
	ID       string   `json:"id"`
	OwnerID  string   `json:"ownerId"`
	Position Position `json:"position"`
	Active   bool     `json:"active"`
}

// Создание нового робота
func NewRobot(id string, ownerID string, level int, position int) Robot {
	return Robot{
		ID:      id,
		OwnerID: ownerID,
		Position: Position{
			Level:    level,
			Position: position,
		},
		Active: true,
	}
}

// Перемещение робота по заранее определённому маршруту
func (r *Robot) Move(path []Position, diceResult int) error {
	if !r.Active {
		return errors.New("робот неактивен")
	}

	if diceResult < 1 || diceResult > 6 {
		return errors.New("некорректный результат кубика")
	}

	if len(path) == 0 {
		return errors.New("маршрут не задан")
	}

	// Ищем текущее положение робота в маршруте
	currentIndex := -1

	for i, position := range path {
		if position == r.Position {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return errors.New("робот не найден на маршруте")
	}

	// Проверяем, достаточно ли клеток для перемещения
	newIndex := currentIndex + diceResult

	if newIndex >= len(path) {
		return errors.New("недостаточно клеток для перемещения")
	}

	// Проверяем допустимость конечной позиции.
	nextPosition := path[newIndex]

	if nextPosition.Level < 1 || nextPosition.Level > 9 {
		return errors.New("некорректный уровень назначения")
	}

	if nextPosition.Position < 1 || nextPosition.Position > 21 {
		return errors.New("некорректная позиция назначения")
	}

	// Перемещаем робота.
	r.Position = nextPosition
	return nil
}
