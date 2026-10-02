package main

import (
	"fmt"
	"math/rand"
)

func printMenu() {
	fmt.Println("Программа КАЛЬКУЛЯТОР НО КРУЧЕ")
	fmt.Println("1. Решение системы линейных уравнений")
	fmt.Println("0. Выход")
}

func createMatrix() [][]float64 {
	var m, n, choose int

	for {
		fmt.Print("Введите количество уравнений и неизвестных -- ")
		_, err := fmt.Scanln(&m, &n)
		if err != nil {
			fmt.Println("Ошибка! Введите два числа через пробел")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		if m != n {
			fmt.Println("Система не имеет одного решения")
			continue
		}
		break
	}

	matrix := make([][]float64, m)
	for i := range matrix {
		matrix[i] = make([]float64, n)
	}

	for {
		fmt.Println("1. Ввести с клавиатуры")
		fmt.Println("2. Заполнить случайно [-10; 10]")
		fmt.Print("Выберите -- ")
		_, err := fmt.Scanln(&choose)
		if err != nil {
			fmt.Println("Ошибка! Введите число")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		if choose != 1 && choose != 2 {
			fmt.Println("Ошибка! Введите 1 или 2")
			continue
		}
		break
	}

	switch choose {
	case 1:
		for i := 0; i < m; i++ {
			fmt.Printf("Строка %d:\n", i+1)
			for j := 0; j < n; j++ {
				var num float64
				for {
					fmt.Printf("a[%d][%d] = ", i+1, j+1)
					_, err := fmt.Scanln(&num)
					if err != nil {
						fmt.Println("Ошибка! Введите число")
						var discard string
						fmt.Scanln(&discard)
						continue
					}
					matrix[i][j] = num
					break
				}
			}
		}
	case 2:
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				matrix[i][j] = float64(rand.Intn(21) - 10)
			}
		}
	}

	return matrix
}

func gaussMetod() {
	m := len(matrix)
	n := len(matrix[0])
	matrix := createMatrix()
	fmt.Println("Матрица:")
	for _, row := range matrix {
		fmt.Printf("%4.4v\n", row)
	}
	for i := 0; i < m; m++ {
		for j := 0; j < n; n++ {
			matrix[i][j] := matrix[i][j] / matrix[i][i]
		}
	}
}

func main() {
	var choose int
	for {
		printMenu()
		_, err := fmt.Scanln(&choose)
		if err != nil {
			fmt.Println("Ошибка! Введите число")
			var discard string
			fmt.Scanln(&discard)
			continue
		}
		switch choose {
		case 0:
			fmt.Println("Выход")
			return
		case 1:
			gaussMetod()
		}
	}
}
