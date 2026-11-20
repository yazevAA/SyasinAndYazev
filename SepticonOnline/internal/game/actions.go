package game

import (
	"errors"

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

	if level < 1 || level > 8 ||
		position < 1 || position > 21 {
		return errors.New("некорректная позиция робота")
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
func (g *Game) StartProduction(
	playerID string,
	taskID string,
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
