package pricemod

import "fmt"

// Считает сумму товаров и применяет прайсмодифаер
func (o Order) String() string {
	var total float64
	for _, item := range o.Items {
		total += item.Price
	}
	if o.Modifier != nil {
		total = o.Modifier(total)
	}
	return fmt.Sprintf("Order{%d items, total=%.2f}", len(o.Items), total)
}
