package board

import (
	"errors"
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()

	if len(b.Levels) != UndergroundLevels {
		t.Fatalf("ожидалось %d уровней, получено %d",
			UndergroundLevels, len(b.Levels))
	}

	if b.Surface.Level != LevelCount {
		t.Fatalf("ожидался уровень поверхности %d", LevelCount)
	}

	for _, level := range b.Levels {
		if len(level.Compartments) != CompartmentsPerLevel {
			t.Errorf("уровень %d содержит неправильное количество отсеков",
				level.Number)
		}

		for _, compartment := range level.Compartments {
			if compartment.Type != CompartmentUnknown {
				t.Errorf("отсек %s должен иметь тип unknown",
					compartment.ID)
			}

			if compartment.Damaged {
				t.Errorf("новый отсек %s не должен быть повреждён",
					compartment.ID)
			}
		}
	}
}

func TestSetCompartmentType(t *testing.T) {
	b := NewBoard()

	err := b.SetCompartmentType("L1-C1", CompartmentProduction)
	if err != nil {
		t.Fatalf("не удалось настроить отсек: %v", err)
	}

	compartment, err := b.FindCompartment("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	if compartment.Type != CompartmentProduction {
		t.Errorf("ожидался тип production, получен %s", compartment.Type)
	}
}

func TestSetInvalidCompartmentType(t *testing.T) {
	b := NewBoard()

	err := b.SetCompartmentType("L1-C1", CompartmentType("invalid"))

	if !errors.Is(err, ErrInvalidCompartmentType) {
		t.Errorf("ожидалась ошибка ErrInvalidCompartmentType, получено %v", err)
	}
}

func TestFindCompartmentNotFound(t *testing.T) {
	b := NewBoard()

	_, err := b.FindCompartment("L99-C99")

	if !errors.Is(err, ErrCompartmentNotFound) {
		t.Errorf("ожидалась ошибка ErrCompartmentNotFound, получено %v", err)
	}
}

func TestCompartmentOperational(t *testing.T) {
	b := NewBoard()

	// Не настроенный отсек недоступен.
	operational, err := b.IsCompartmentOperational("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	if operational {
		t.Error("отсек unknown не должен быть доступен")
	}

	// Настроенный и неповреждённый отсек доступен.
	err = b.SetCompartmentType("L1-C1", CompartmentEnergy)
	if err != nil {
		t.Fatal(err)
	}

	operational, err = b.IsCompartmentOperational("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	if !operational {
		t.Error("исправный настроенный отсек должен быть доступен")
	}

	// Повреждённый отсек недоступен.
	err = b.DamageCompartment("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	operational, err = b.IsCompartmentOperational("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	if operational {
		t.Error("повреждённый отсек не должен быть доступен")
	}
}

func TestRepairCompartment(t *testing.T) {
	b := NewBoard()

	b.SetCompartmentType("L1-C1", CompartmentProduction)
	b.DamageCompartment("L1-C1")

	err := b.RepairCompartment("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	operational, err := b.IsCompartmentOperational("L1-C1")
	if err != nil {
		t.Fatal(err)
	}

	if !operational {
		t.Error("отсек должен стать доступным после ремонта")
	}
}

func TestValidateBoard(t *testing.T) {
	b := NewBoard()

	if err := b.Validate(); err != nil {
		t.Fatalf("новая станция должна быть корректной: %v", err)
	}

	b.Levels[0].Compartments[0].Type = CompartmentType("invalid")

	if err := b.Validate(); !errors.Is(err, ErrInvalidCompartmentType) {
		t.Errorf("ожидалась ошибка недопустимого типа отсека, получено %v", err)
	}
}
