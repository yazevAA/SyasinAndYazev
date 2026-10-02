package main

import (
	"fmt"
	"math"
)

func printMenu() {
	fmt.Println("\nПрограмма КАЛЬКУЛЯТОР может вычислять:")
	fmt.Println("1.Вычислять площадь треугольника по трем сторонам")
	fmt.Println("2.Вычислить десятичный логарифм от числа")
	fmt.Println("3.Найти квадрат числа")
	fmt.Println("4.Найти Косинус и Синус угла, заданного в радианах")
	fmt.Println("5.Кубический корень числа")
	fmt.Println("0.Закрыть программу")
	fmt.Print("Выберите пункт -- ")
}

func triangleSqr() {
	fmt.Println()
	fmt.Println("Расчет площади треугольника по трем сторонам")
	var a, b, c float64
	for {
		fmt.Print("Введите стороны a, b, c: ")
		_, err := fmt.Scanln(&a, &b, &c)
		if err != nil {
			fmt.Println("Ошибка! Введите три числа через пробел")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		if a <= 0 || b <= 0 || c <= 0 || a+b <= c || a+c <= b || b+c <= a {
			fmt.Println("Не выполняется условие существования треугольника")
			continue
		}
		break
	}

	p := (a + b + c) / 2
	sqr := math.Sqrt(p * (p - a) * (p - b) * (p - c))
	fmt.Printf("Площадь треугольника = %.4f\n", sqr)
}

func lg() {
	fmt.Println()
	fmt.Println("Расчет десятичного логарифма от числа")
	var x, lg float64
	for {
		fmt.Print("Введите число -- ")
		_, err := fmt.Scanln(&x)
		if err != nil || x <= 0 {
			fmt.Println("Ошибка! Введите положительное число")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		break
	}

	lg = math.Log10(x)
	fmt.Printf("lg = %.4f\n", lg)
}

func sqr() {
	fmt.Println()
	fmt.Println("Расчет квадрата числа")
	var x, sqr float64
	fmt.Print("Введите число -- ")
	fmt.Scanln(&x)

	sqr = math.Pow(x, 2)
	fmt.Printf("Квадрат числа = %.4f\n", sqr)
}

func cosAndSin() {
	fmt.Println()
	fmt.Println("Расчет синуса и косинуса угла")
	var x, sin, cos float64

	fmt.Print("Введите значение угла -- ")
	fmt.Scanln(&x)

	sin = math.Sin(x)
	cos = math.Cos(x)

	fmt.Printf("Синус угла = %.4f  Косинус угла = %.4f\n", sin, cos)
}

func cubicheskiRoot() {
	fmt.Println()
	fmt.Println("Расчет кубического корня из числа")
	var x, cbrt float64
	for {
		fmt.Print("Введите число -- ")
		_, err := fmt.Scanln(&x)
		if err != nil || x < 0 {
			fmt.Println("Ошибка! Введите неотрицательное число")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		break
	}

	cbrt = math.Cbrt(x)
	fmt.Printf("Кубический корень числа = %.4f\n", cbrt)
}

func main() {
	for {
		printMenu()

		var choose uint8
		_, err := fmt.Scanln(&choose)
		if err != nil {
			fmt.Println("Ошибка! Введите целое число от 0 до 5")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		if choose > 5 {
			fmt.Println("Ошибка! Введите целое число от 0 до 5")
			continue
		}

		switch choose {
		case 0:
			fmt.Println("Закрытие")
			return
		case 1:
			triangleSqr()
		case 2:
			lg()
		case 3:
			sqr()
		case 4:
			cosAndSin()
		case 5:
			cubicheskiRoot()
		}
	}
}
