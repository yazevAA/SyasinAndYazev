package pricing

import (
	"fmt"
	"math"
)

// PriceModifier — тип-функция, преобразующая цену.
type PriceModifier func(price float64) float64

// WithDiscount возвращает модификатор скидки percent процентов.
func WithDiscount(percent float64) PriceModifier {
	if percent < 0 || percent > 100 {
		panic(fmt.Sprintf("pricing: недопустимая скидка: %v", percent))
	}
	pct := percent
	return func(price float64) float64 {
		return price * (1 - pct/100)
	}
}

// WithTax возвращает модификатор налога rate процентов.
func WithTax(rate float64) PriceModifier {
	if rate < 0 {
		panic(fmt.Sprintf("pricing: недопустимая ставка налога: %v", rate))
	}
	r := rate
	return func(price float64) float64 {
		return price * (1 + r/100)
	}
}

// WithRounding возвращает модификатор округления до precision знаков.
func WithRounding(precision int) PriceModifier {
	if precision < 0 || precision > 8 {
		panic(fmt.Sprintf("pricing: недопустимая точность: %d", precision))
	}
	factor := math.Pow(10, float64(precision))
	return func(price float64) float64 {
		return math.Round(price*factor) / factor
	}
}

// WithMinimumPrice поднимает цену до minPrice.
func WithMinimumPrice(minPrice float64) PriceModifier {
	if minPrice < 0 {
		panic(fmt.Sprintf("pricing: отрицательная минимальная цена: %v", minPrice))
	}
	min := minPrice
	return func(price float64) float64 {
		if price < min {
			return min
		}
		return price
	}
}

// WithMaximumPrice ограничивает цену сверху.
func WithMaximumPrice(maxPrice float64) PriceModifier {
	if maxPrice < 0 {
		panic(fmt.Sprintf("pricing: отрицательная максимальная цена: %v", maxPrice))
	}
	max := maxPrice
	return func(price float64) float64 {
		if price > max {
			return max
		}
		return price
	}
}

// WithFixedDiscount вычитает фиксированную сумму; цена не уходит в минус.
func WithFixedDiscount(amount float64) PriceModifier {
	if amount < 0 {
		panic(fmt.Sprintf("pricing: отрицательная фиксированная скидка: %v", amount))
	}
	a := amount
	return func(price float64) float64 {
		res := price - a
		if res < 0 {
			return 0
		}
		return res
	}
}
