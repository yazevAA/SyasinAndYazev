package pricing

import "fmt"

// Product описывает товар.
type Product struct {
	ID    string
	Name  string
	Price float64
}

// String реализует fmt.Stringer.
func (p Product) String() string {
	return fmt.Sprintf("Товар{артикул=%s название=%q цена=%.2f}", p.ID, p.Name, p.Price)
}

// Total возвращает стоимость quantity единиц.
func (p Product) Total(quantity int) float64 {
	if quantity <= 0 {
		return 0
	}
	return p.Price * float64(quantity)
}
