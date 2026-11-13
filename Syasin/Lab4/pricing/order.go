package pricing

import "fmt"

// OrderItem — позиция заказа.
type OrderItem struct {
	Product  Product
	Quantity int
}

// Subtotal возвращает базовую стоимость позиции.
func (it OrderItem) Subtotal() float64 {
	return it.Product.Total(it.Quantity)
}

// Order описывает заказ.
type Order struct {
	ID    string
	Items []OrderItem
}

// Subtotal суммирует базовые цены позиций.
func (o Order) Subtotal() float64 {
	var sum float64
	for _, it := range o.Items {
		sum += it.Subtotal()
	}
	return sum
}

// String реализует fmt.Stringer.
func (o Order) String() string {
	return fmt.Sprintf("Заказ{номер=%s позиций=%d сумма=%.2f}", o.ID, len(o.Items), o.Subtotal())
}

// FinalPrice применяет модификатор к Subtotal.
func (o Order) FinalPrice(mod PriceModifier) float64 {
	if mod == nil {
		return o.Subtotal()
	}
	return mod(o.Subtotal())
}
