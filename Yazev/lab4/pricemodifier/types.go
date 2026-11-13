package pricemod

// Тут товар и его цена
type Product struct {
	Name  string
	Price float64
}

// сам "заказ" список товаров и применение скидки
type Order struct {
	Items    []Product
	Modifier PriceModifier //как посчитать новую цену
}
