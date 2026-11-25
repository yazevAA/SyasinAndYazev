package game

import (
	"errors"

	"septikon-online/internal/board"
	"septikon-online/internal/resources"
	"septikon-online/internal/robot"
	"septikon-online/internal/rocket"
)

// Бросок кубика
func (g *Game) RollDice(playerID string) (int, error) {
	if err := g.CheckTurn(playerID); err != nil {
		return 0, err
	}

	return g.TurnManager.RollDice(playerID)
}

// Добавление робота игроку
func (g *Game) AddRobot(
	playerID string,
	robotID string,
	level int,
	position int,
) error {
	p, err := g.FindPlayer(playerID)
	if err != nil {
		return err
	}

	if robotID == "" {
		return errors.New("ID робота не может быть пустым")
	}

	// Проверяем допустимость положения робота.
	if level < 1 || level > board.LevelCount {
		return errors.New("недопустимый уровень")
	}

	// На подземных уровнях доступно 21 позиция.
	// На поверхности пока используем ту же систему колонок.
	if position < 1 || position > board.CompartmentsPerLevel {
		return errors.New("недопустимая позиция робота")
	}

	// Проверяем наличие поверхности.
	if level == board.LevelCount && !p.Station.Surface.Present {
		return errors.New("поверхность станции недоступна")
	}

	for _, existing := range p.Robots {
		if existing.ID == robotID {
			return errors.New("робот с таким ID уже существует")
		}
	}

	p.Robots = append(p.Robots,
		robot.NewRobot(robotID, playerID, level, position),
	)

	return nil
}

// FindGunner находит активного робота игрока на поверхности станции.
func (g *Game) FindGunner(
	playerID string,
	robotID string,
) (*robot.Robot, error) {
	player, err := g.FindPlayer(playerID)
	if err != nil {
		return nil, err
	}

	for i := range player.Robots {
		r := &player.Robots[i]

		if r.ID == robotID {
			if !r.Active {
				return nil, errors.New("робот неактивен")
			}

			if r.Position.Level != 9 {
				return nil, errors.New("робот не находится на поверхности")
			}

			return r, nil
		}
	}

	return nil, errors.New("робот не найден")
}

// Перемещение робота
func (g *Game) MoveRobot(
	playerID string,
	robotID string,
	path []robot.Position,
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

	p, err := g.FindPlayer(playerID)
	if err != nil {
		return err
	}

	for i := range p.Robots {
		if p.Robots[i].ID == robotID {
			if err := p.Robots[i].Move(
				path,
				g.TurnManager.Turn.DiceResult,
			); err != nil {
				return err
			}

			return g.TurnManager.UseAction(playerID)
		}
	}

	return errors.New("робот не найден")
}

// Запуск производства ресурсов
// Запуск производства ресурсов в указанном отсеке.
func (g *Game) StartProduction(
	playerID string,
	taskID string,
	compartmentID string,
	resource resources.ResourceType,
	amount int,
	requiredTurns int,
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

	p, err := g.FindPlayer(playerID)
	if err != nil {
		return err
	}

	if taskID == "" {
		return errors.New("ID задачи не может быть пустым")
	}

	if compartmentID == "" {
		return errors.New("ID производственного отсека не может быть пустым")
	}

	// Проверяем, существует ли отсек на станции игрока.
	compartment, err := p.Station.FindCompartment(compartmentID)
	if err != nil {
		return err
	}

	// Производство разрешено только в производственном отсеке.
	if compartment.Type != board.CompartmentProduction {
		return errors.New("отсек не предназначен для производства")
	}

	// Повреждённый отсек использовать нельзя.
	if compartment.Damaged {
		return errors.New("производственный отсек повреждён")
	}

	// Проверяем уникальность задачи.
	for _, task := range p.ProductionTasks {
		if task.ID == taskID {
			return errors.New("задача с таким ID уже существует")
		}
	}

	task, err := resources.NewProductionTask(
		taskID,
		resource,
		amount,
		requiredTurns,
	)
	if err != nil {
		return err
	}

	task.CompartmentID = compartmentID

	p.ProductionTasks = append(p.ProductionTasks, task)

	return g.TurnManager.UseAction(playerID)
}

// Получение готовых ресурсов
func (g *Game) CollectProduction(
	playerID string,
	taskID string,
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

	p, err := g.FindPlayer(playerID)
	if err != nil {
		return err
	}

	for i := range p.ProductionTasks {
		task := &p.ProductionTasks[i]

		if task.ID == taskID {
			if err := task.Collect(&p.Inventory); err != nil {
				return err
			}

			return g.TurnManager.UseAction(playerID)
		}
	}

	return errors.New("задача добычи не найдена")
}

// Запуск ракеты
func (g *Game) LaunchRocket(
	playerID string,
	rocketID string,
	gunnerID string,
	targetStation string,
	distance int,
	speed int,
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

	if rocketID == "" {
		return errors.New("ID ракеты не может быть пустым")
	}

	for _, existing := range g.Rockets {
		if existing.ID == rocketID {
			return errors.New("ракета с таким ID уже существует")
		}
	}

	// Проверяем наличие активного робота-стрелка на поверхности.
	if _, err := g.FindGunner(playerID, gunnerID); err != nil {
		return err
	}

	targetPlayer, err := g.FindPlayer(targetStation)
	if err != nil {
		return errors.New("целевая станция не найдена")
	}

	if targetPlayer.ID == playerID {
		return errors.New("нельзя запускать ракету на свою станцию")
	}

	newRocket, err := rocket.NewRocket(
		rocketID,
		playerID,
		playerID,
		targetStation,
		distance,
		speed,
	)

	if err != nil {
		return err
	}

	if err := newRocket.Launch(); err != nil {
		return err
	}

	g.Rockets = append(g.Rockets, newRocket)

	return g.TurnManager.UseAction(playerID)
}

// Обработка игровых процессов в конце хода

func (g *Game) AdvanceProcesses() {
	// Продвижение ракет
	for i := range g.Rockets {
		if g.Rockets[i].Status == rocket.RocketFlying {
			_ = g.Rockets[i].AdvanceTurn()
		}
	}

	// Продвижение добычи ресурсов
	for i := range g.Players {
		for j := range g.Players[i].ProductionTasks {
			task := &g.Players[i].ProductionTasks[j]

			if !task.Completed {
				_ = task.AdvanceTurn()
			}
		}
	}

	// Продвижение ремонта отсеков
	g.AdvanceRepairs()
}

// Завершение хода
func (g *Game) EndTurn(playerID string) error {
	if err := g.CheckTurn(playerID); err != nil {
		return err
	}

	if !g.TurnManager.Turn.DiceRolled {
		return errors.New("сначала необходимо бросить кубик")
	}

	if !g.TurnManager.Turn.ActionUsed {
		return errors.New("необходимо выполнить основное действие")
	}

	// Продвигаем игровые процессы
	g.AdvanceProcesses()

	// Передаём ход следующему игроку
	return g.TurnManager.EndTurn(playerID)
}
