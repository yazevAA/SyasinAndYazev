package main

import (
	"fmt"
	"math"
	"strings"
)

// Структура для хранения результатов вычислений
type Result struct {
	Method     string
	Approx     float64
	Exact      float64
	Error      float64
	Iterations int
}

// String() задаёт красивый формат вывода результата
func (r Result) String() string {
	return fmt.Sprintf("%-23s | %14.8f | %14.8f | %12.3e | %8d",
		r.Method, r.Approx, r.Exact, r.Error, r.Iterations)
}

// Метод Ньютона для вычисления квадратного корня
func newtonSqrt(a, eps float64) Result {
	exact := math.Sqrt(a)

	x := a
	if x < 1 {
		x = 1
	}

	iterations := 0
	const maxIterations = 1_000_000

	for iterations < maxIterations {
		next := (x + a/x) / 2
		iterations++

		difference := math.Abs(next - x)
		x = next

		// Останавливаем вычисления при достижении точности
		if difference <= eps {
			break
		}
	}

	return Result{
		Method:     "Метод Ньютона",
		Approx:     x,
		Exact:      exact,
		Error:      math.Abs(x - exact),
		Iterations: iterations,
	}
}

// Метод деления пополам
func bisectionSqrt(a, eps float64) Result {
	exact := math.Sqrt(a)

	left, right := 0.0, math.Max(1, a)
	iterations := 0

	const maxIterations = 1_000_000

	for iterations < maxIterations {
		middle := (left + right) / 2
		iterations++

		// Проверяем достижение необходимой точности
		if (right-left)/2 <= eps ||
			middle == left || middle == right {
			break
		}

		if middle*middle == a {
			left, right = middle, middle
			break
		}

		if middle*middle < a {
			left = middle
			continue
		}

		right = middle
	}

	approx := (left + right) / 2

	return Result{
		Method:     "Метод деления пополам",
		Approx:     approx,
		Exact:      exact,
		Error:      math.Abs(approx - exact),
		Iterations: iterations,
	}
}

// Вывод таблицы результатов
func PrintTable(results []Result) {
	fmt.Println("\nСРАВНЕНИЕ РЕЗУЛЬТАТОВ")
	fmt.Println(strings.Repeat("-", 91))

	fmt.Printf("%-23s | %14s | %14s | %12s | %8s\n",
		"Метод", "Приближённое", "Точное",
		"Погрешность", "Итерации")

	fmt.Println(strings.Repeat("-", 91))

	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println(strings.Repeat("-", 91))
}

func main() {
	var a, eps float64

	fmt.Println("ЛАБОРАТОРИЯ ЧИСЛЕННЫХ МЕТОДОВ")
	fmt.Println("Сравнение способов вычисления квадратного корня")

	// Ввод исходного числа
	fmt.Print("Введите положительное число a: ")

	if _, err := fmt.Scan(&a); err != nil {
		fmt.Println("Ошибка: необходимо ввести число.")
		return
	}

	if a <= 0 || math.IsNaN(a) || math.IsInf(a, 0) {
		fmt.Println("Ошибка: число a должно быть положительным и конечным.")
		return
	}

	// Ввод точности
	fmt.Print("Введите точность ε (например, 0.000001): ")

	if _, err := fmt.Scan(&eps); err != nil {
		fmt.Println("Ошибка: необходимо ввести число.")
		return
	}

	if eps <= 0 || math.IsNaN(eps) || math.IsInf(eps, 0) {
		fmt.Println("Ошибка: точность ε должна быть положительной.")
		return
	}

	// Выполнение вычислений двумя методами
	results := []Result{
		newtonSqrt(a, eps),
		bisectionSqrt(a, eps),
	}

	// Вывод результатов
	PrintTable(results)

	// Определение более быстрого метода
	if results[0].Iterations < results[1].Iterations {
		fmt.Println("Быстрее по количеству итераций: метод Ньютона.")
	} else if results[1].Iterations < results[0].Iterations {
		fmt.Println("Быстрее по количеству итераций: метод деления пополам.")
	} else {
		fmt.Println("Оба метода выполнили одинаковое количество итераций.")
	}
}
