package rocket

import "errors"

// Состояние ракеты
type RocketStatus string

const (
	RocketReady   RocketStatus = "ready"
	RocketFlying  RocketStatus = "flying"
	RocketArrived RocketStatus = "arrived"
)

// Игровая ракета
type Rocket struct {
	ID            string       `json:"id"`
	OwnerID       string       `json:"ownerId"`
	SourceStation string       `json:"sourceStation"`
	TargetStation string       `json:"targetStation"`
	Position      int          `json:"position"`
	TotalDistance int          `json:"totalDistance"`
	Speed         int          `json:"speed"`
	Status        RocketStatus `json:"status"`
}

// Создание ракеты
func NewRocket(
	id string,
	ownerID string,
	source string,
	target string,
	distance int,
	speed int,
) (Rocket, error) {

	if distance <= 0 || speed <= 0 {
		return Rocket{}, errors.New("некорректные параметры ракеты")
	}

	if source == target {
		return Rocket{}, errors.New("станции должны различаться")
	}

	return Rocket{
		ID:            id,
		OwnerID:       ownerID,
		SourceStation: source,
		TargetStation: target,
		Position:      0,
		TotalDistance: distance,
		Speed:         speed,
		Status:        RocketReady,
	}, nil
}

// Запуск ракеты
func (r *Rocket) Launch() error {
	if r.Status != RocketReady {
		return errors.New("ракета не готова к запуску")
	}

	r.Status = RocketFlying

	return nil
}

// Перемещение ракеты на один игровой ход
func (r *Rocket) AdvanceTurn() error {
	if r.Status != RocketFlying {
		return errors.New("ракета не находится в полёте")
	}

	r.Position += r.Speed

	// Проверяем достижение цели
	if r.Position >= r.TotalDistance {
		r.Position = r.TotalDistance
		r.Status = RocketArrived
	}

	return nil
}
