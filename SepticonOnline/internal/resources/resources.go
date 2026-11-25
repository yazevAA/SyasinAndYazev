package resources

import "errors"

// Типы игровых ресурсов
type ResourceType string

const (
	Metal   ResourceType = "metal"
	Uranium ResourceType = "uranium"
	Oxygen  ResourceType = "oxygen"
	Biomass ResourceType = "biomass"
	Energy  ResourceType = "energy"
)

// Запасы ресурсов игрока
type Inventory struct {
	Metal   int `json:"metal"`
	Uranium int `json:"uranium"`
	Oxygen  int `json:"oxygen"`
	Biomass int `json:"biomass"`
	Energy  int `json:"energy"`
}

// Добавление ресурсов
func (i *Inventory) Add(resource ResourceType, amount int) error {
	if amount <= 0 {
		return errors.New("количество должно быть положительным")
	}

	switch resource {
	case Metal:
		i.Metal += amount
	case Uranium:
		i.Uranium += amount
	case Oxygen:
		i.Oxygen += amount
	case Biomass:
		i.Biomass += amount
	case Energy:
		i.Energy += amount
	default:
		return errors.New("неизвестный тип ресурса")
	}

	return nil
}

// Проверка наличия ресурсов
func (i Inventory) Has(resource ResourceType, amount int) bool {
	if amount < 0 {
		return false
	}

	switch resource {
	case Metal:
		return i.Metal >= amount
	case Uranium:
		return i.Uranium >= amount
	case Oxygen:
		return i.Oxygen >= amount
	case Biomass:
		return i.Biomass >= amount
	case Energy:
		return i.Energy >= amount
	default:
		return false
	}
}

// Затрата ресурсов
func (i *Inventory) Spend(resource ResourceType, amount int) error {
	if amount <= 0 {
		return errors.New("количество должно быть положительным")
	}

	if !i.Has(resource, amount) {
		return errors.New("недостаточно ресурсов")
	}

	switch resource {
	case Metal:
		i.Metal -= amount
	case Uranium:
		i.Uranium -= amount
	case Oxygen:
		i.Oxygen -= amount
	case Biomass:
		i.Biomass -= amount
	case Energy:
		i.Energy -= amount
	default:
		return errors.New("неизвестный тип ресурса")
	}

	return nil
}

// Задача по добыче ресурсов
// Задача по добыче ресурсов.
type ProductionTask struct {
	ID             string       `json:"id"`
	CompartmentID  string       `json:"compartmentId"`
	Resource       ResourceType `json:"resource"`
	Amount         int          `json:"amount"`
	RequiredTurns  int          `json:"requiredTurns"`
	RemainingTurns int          `json:"remainingTurns"`
	Completed      bool         `json:"completed"`
	Collected      bool         `json:"collected"`
}

// Создание задачи добычи
func NewProductionTask(
	id string,
	resource ResourceType,
	amount int,
	requiredTurns int,
) (ProductionTask, error) {

	if amount <= 0 || requiredTurns <= 0 {
		return ProductionTask{}, errors.New("некорректные параметры добычи")
	}

	switch resource {
	case Metal, Uranium, Oxygen, Biomass, Energy:
	default:
		return ProductionTask{}, errors.New("неизвестный тип ресурса")
	}

	return ProductionTask{
		ID:             id,
		Resource:       resource,
		Amount:         amount,
		RequiredTurns:  requiredTurns,
		RemainingTurns: requiredTurns,
		Completed:      false,
	}, nil
}

// Продвижение задачи на один ход
func (t *ProductionTask) AdvanceTurn() error {
	if t.Completed {
		return errors.New("добыча уже завершена")
	}

	t.RemainingTurns--

	if t.RemainingTurns <= 0 {
		t.RemainingTurns = 0
		t.Completed = true
	}

	return nil
}

// Получение добытых ресурсов
// Получение добытых ресурсов.
func (t *ProductionTask) Collect(inventory *Inventory) error {
	if !t.Completed {
		return errors.New("добыча еще не завершена")
	}

	if t.Collected {
		return errors.New("ресурсы уже получены")
	}

	if inventory == nil {
		return errors.New("инвентарь не существует")
	}

	if err := inventory.Add(t.Resource, t.Amount); err != nil {
		return err
	}

	// Отмечаем, что ресурсы уже получены.
	// Completed остаётся true, Amount сохраняет исходное значение.
	t.Collected = true

	return nil
}
