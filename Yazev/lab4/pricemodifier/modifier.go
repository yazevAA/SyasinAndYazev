package pricemod

// функция изменяющая цену
type PriceModifier func(float64) float64

// применение скидки
func WithDiscount(percent float64) PriceModifier {
	factor := 1 - percent/100
	return func(price float64) float64 {
		return price * factor
	}
}

// применение налога
func WithTax(rate float64) PriceModifier {
	factor := 1 + rate/100
	return func(price float64) float64 {
		return price * factor
	}
}

// объединяет модификаторы в один, применяя их по порядку
func CombineModifiers(modifiers ...PriceModifier) PriceModifier {
	return func(price float64) float64 {
		for _, m := range modifiers {
			price = m(price)
		}
		return price
	}
}
