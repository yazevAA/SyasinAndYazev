package main

import (
	"fmt"

	"pricemod"
)

func main() {
	order := pricemod.Order{
		Items: []pricemod.Product{
			{"Ледяная скорбь", 17000},
			{"Испепелитель", 20000},
			{"Копье Фавония", 500},
		},
	}

	//  обычные цены
	fmt.Println("         Обычные цены")
	printTable(order.Items, nil)

	// запрос скидки и налога
	var discount, tax float64
	fmt.Print("Введите скидку в %: ")
	fmt.Scanln(&discount)
	fmt.Print("Введите налог в %: ")
	fmt.Scanln(&tax)

	//  модификатор из введённых чисел
	order.Modifier = pricemod.CombineModifiers(
		pricemod.WithDiscount(discount),
		pricemod.WithTax(tax),
	)

	//     со скидкой и налогом
	fmt.Printf("\n       Со скидкой %.2f%% и налогом %.2f%%\n", discount, tax)
	printTable(order.Items, order.Modifier)
}

func printTable(items []pricemod.Product, mod pricemod.PriceModifier) {
	if mod == nil {
		fmt.Printf("%-20s %10s\n", "Товар", "Цена")
	} else {
		fmt.Printf("%-20s %10s %12s\n", "Товар", "Цена", "Со скидкой")
	}
	fmt.Println("------------------------------------------")

	var total, totalFinal float64
	for _, it := range items {
		total += it.Price
		if mod == nil {
			fmt.Printf("%-20s %10.2f\n", it.Name, it.Price)
		} else {
			final := mod(it.Price)
			totalFinal += final
			fmt.Printf("%-20s %10.2f %12.2f\n", it.Name, it.Price, final)
		}
	}
	fmt.Println("------------------------------------------")
	if mod == nil {
		fmt.Printf("%-20s %10.2f\n", "Итого", total)
	} else {
		fmt.Printf("%-20s %10.2f %12.2f\n", "Итого", total, totalFinal)
	}
}
