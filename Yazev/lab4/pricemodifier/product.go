package pricemod

import "fmt"

// читаемый вывод товара
func (p Product) String() string {
	return fmt.Sprintf("%s: %.2f", p.Name, p.Price)
}
