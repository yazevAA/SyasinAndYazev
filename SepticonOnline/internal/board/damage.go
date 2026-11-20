package board

import "errors"

// Поиск отсека по идентификатору
func (b *Board) FindCompartment(id string) (*Compartment, error) {
	for i := range b.Levels {
		for j := range b.Levels[i].Compartments {
			compartment := &b.Levels[i].Compartments[j]

			if compartment.ID == id {
				return compartment, nil
			}
		}
	}

	return nil, errors.New("отсек не найден")
}

// Повреждение отсека
func (b *Board) DamageCompartment(id string) error {
	compartment, err := b.FindCompartment(id)

	if err != nil {
		return err
	}

	compartment.Damaged = true

	return nil
}

// Восстановление отсека
func (b *Board) RepairCompartment(id string) error {
	compartment, err := b.FindCompartment(id)

	if err != nil {
		return err
	}

	compartment.Damaged = false

	return nil
}
