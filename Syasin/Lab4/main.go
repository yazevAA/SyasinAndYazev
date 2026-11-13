package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lab4/pricing"
)

// --- Константы ---
const (
	TaxRate           = 20.0 // НДС в РФ
	RoundingPrecision = 2    // Округление до копеек
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	fmt.Println("=== Добро пожаловать в TechStore! ===")

	order := pricing.Order{
		ID:    "ORDER-001",
		Items: []pricing.OrderItem{},
	}
	currentDiscount := 0.0

	for {
		fmt.Println("\n--- ГЛАВНОЕ МЕНЮ ---")
		fmt.Println("1. Перейти в каталог")
		fmt.Println("2. Назначить скидку")
		fmt.Println("3. Показать текущую корзину")
		fmt.Println("4. Оформить заказ")
		fmt.Println("5. Выйти")
		fmt.Print("Выберите действие: ")

		choice := readInput()

		switch choice {
		case "1":
			catalogMenu(&order) // Вызываем меню каталога
		case "2":
			setDiscount(&currentDiscount)
		case "3":
			showCart(order, currentDiscount)
		case "4":
			checkout(&order, currentDiscount)
		case "5":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный ввод. Пожалуйста, введите число от 1 до 5.")
		}
	}
}

// readInput читает строку с консоли и убирает лишние пробелы
func readInput() string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// catalogMenu зацикливает процесс добавления товаров
func catalogMenu(order *pricing.Order) {
	for {
		fmt.Println("\n--- КАТАЛОГ ТОВАРОВ ---")
		for i, p := range pricing.Catalog {
			fmt.Printf("%d. %-25s | Цена: %8.2f руб.\n", i+1, p.Name, p.Price)
		}
		fmt.Println("0. Вернуться в главное меню")
		fmt.Print("Выберите номер товара для добавления (или 0 для выхода): ")

		choiceStr := readInput()

		if choiceStr == "0" {
			return // Выходим из каталога обратно в главное меню
		}

		idx, err := strconv.Atoi(choiceStr)
		if err != nil || idx < 1 || idx > len(pricing.Catalog) {
			fmt.Println("Ошибка: неверный номер товара.")
			continue // Возвращаемся к началу цикла каталога
		}

		// Берем товар из каталога
		product := pricing.Catalog[idx-1]

		fmt.Printf("Вы выбрали: %s\n", product.Name)
		fmt.Print("Введите количество (или 0 чтобы отменить): ")
		qtyStr := readInput()

		qty, err := strconv.Atoi(qtyStr)
		if err != nil || qty < 0 {
			fmt.Println("Ошибка: количество должно быть целым неотрицательным числом.")
			continue
		}

		if qty == 0 {
			fmt.Println("Добавление отменено.")
			continue
		}

		// Проверяем, есть ли уже такой товар в корзине
		found := false
		for i, item := range order.Items {
			if item.Product.ID == product.ID {
				order.Items[i].Quantity += qty
				fmt.Printf("✓ Количество товара '%s' увеличено. Теперь в корзине: %d шт.\n", product.Name, order.Items[i].Quantity)
				found = true
				break
			}
		}

		// Если товара еще нет, добавляем новую позицию
		if !found {
			order.Items = append(order.Items, pricing.OrderItem{
				Product:  product,
				Quantity: qty,
			})
			fmt.Printf("✓ Товар '%s' (x%d) добавлен в корзину!\n", product.Name, qty)
		}
	}
}

func setDiscount(currentDiscount *float64) {
	fmt.Printf("\nТекущая скидка: %.2f%%\n", *currentDiscount)
	fmt.Print("Введите новую скидку в процентах (от 0 до 100): ")

	discountStr := readInput()
	discount, err := strconv.ParseFloat(discountStr, 64)
	if err != nil || discount < 0 || discount > 100 {
		fmt.Println("Ошибка: скидка должна быть числом от 0 до 100.")
		return
	}

	*currentDiscount = discount
	fmt.Printf("Скидка %.2f%% успешно установлена.\n", discount)
}

func showCart(order pricing.Order, discount float64) {
	if len(order.Items) == 0 {
		fmt.Println("\nКорзина пуста.")
		return
	}

	fmt.Println("\n--- ВАША КОРЗИНА ---")
	for i, item := range order.Items {
		fmt.Printf("%d. %-25s | Цена: %8.2f | Кол-во: %2d | Сумма: %8.2f\n",
			i+1, item.Product.Name, item.Product.Price, item.Quantity, item.Subtotal())
	}
	fmt.Printf("\nИтого позиций: %d\n", len(order.Items))
	fmt.Printf("Базовая сумма (без НДС и скидки): %.2f руб.\n", order.Subtotal())
	fmt.Printf("Текущая скидка: %.2f%%\n", discount)
}

func checkout(order *pricing.Order, discount float64) {
	if len(order.Items) == 0 {
		fmt.Println("\nКорзина пуста. Нечего оформлять.")
		return
	}

	fmt.Println("\n--- ОФОРМЛЕНИЕ ЗАКАЗА ---")

	modifiers := []pricing.PriceModifier{}

	if discount > 0 {
		modifiers = append(modifiers, pricing.WithDiscount(discount))
	}

	modifiers = append(modifiers, pricing.WithTax(TaxRate))
	modifiers = append(modifiers, pricing.WithRounding(RoundingPrecision))

	finalPriceMod := pricing.CombineModifiers(modifiers...)

	subtotal := order.Subtotal()

	// Пошаговый расчет для наглядности
	fmt.Printf("1. Базовая сумма:             %8.2f руб.\n", subtotal)

	currentSum := subtotal
	if discount > 0 {
		currentSum = pricing.WithDiscount(discount)(currentSum)
		fmt.Printf("2. После скидки %.1f%%:         %8.2f руб.\n", discount, currentSum)
	} else {
		fmt.Printf("2. Без скидки:                %8.2f руб.\n", currentSum)
	}

	// НДС считается от суммы ПОСЛЕ скидки
	currentSum = pricing.WithTax(TaxRate)(currentSum)
	fmt.Printf("3. С НДС %v%%:                  %8.2f руб.\n", TaxRate, currentSum)

	currentSum = pricing.WithRounding(RoundingPrecision)(currentSum)
	fmt.Printf("4. С округлением до %d знаков: %8.2f руб.\n", RoundingPrecision, currentSum)

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("ИТОГО К ОПЛАТЕ:               %8.2f руб.\n", order.FinalPrice(finalPriceMod))
	fmt.Println(strings.Repeat("-", 50))

	// Очищаем корзину
	order.Items = []pricing.OrderItem{}
	fmt.Println("\nЗаказ успешно оформлен! Корзина очищена.")
}
