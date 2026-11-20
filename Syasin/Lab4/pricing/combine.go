package pricing

// CombineModifiers объединяет произвольное число модификаторов
// в один, применяя их последовательно слева направо.
func CombineModifiers(modifiers ...PriceModifier) PriceModifier {
	switch len(modifiers) {
	case 0:
		return func(p float64) float64 { return p }
	case 1:
		return modifiers[0]
	}
	mods := append([]PriceModifier(nil), modifiers...)
	return func(price float64) float64 {
		for _, m := range mods {
			price = m(price)
		}
		return price
	}
}
