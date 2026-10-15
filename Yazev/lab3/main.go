package main

import (
	"fmt"
	"math"
)

func sinTaylor(x, eps float64, maxIter int) (float64, int) {
	x = math.Mod(x, 2*math.Pi)
	if x > math.Pi {
		x -= 2 * math.Pi
	}
	if x < -math.Pi {
		x += 2 * math.Pi
	}

	term := x
	result := x
	n := 1

	for i := 1; i <= maxIter; i++ {
		term *= -x * x / float64((2*n)*(2*n+1))
		result += term
		n++

		if math.Abs(term) < eps {
			return result, i
		}
	}
	return result, maxIter
}

func cosTaylor(x, eps float64, maxIter int) (float64, int) {
	x = math.Mod(x, 2*math.Pi)
	if x > math.Pi {
		x -= 2 * math.Pi
	}
	if x < -math.Pi {
		x += 2 * math.Pi
	}

	term := 1.0
	result := 1.0
	n := 1

	for i := 1; i <= maxIter; i++ {
		term *= -x * x / float64((2*n-1)*(2*n))
		result += term
		n++

		if math.Abs(term) < eps {
			return result, i
		}
	}
	return result, maxIter
}

func piNilakantha(eps float64, maxIter int) (float64, int) {
	total := 3.0
	sign := 1.0
	n := 2

	for i := 1; i <= maxIter; i++ {
		term := 4.0 / float64(n*(n+1)*(n+2))
		total += sign * term
		sign = -sign
		n += 2

		if term > eps {
			continue
		}

		return total, i
	}
	return total, maxIter
}

func eTaylor(eps float64, maxIter int) (float64, int) {
	term := 1.0
	total := 1.0
	n := 1

	for {
		if n >= maxIter {
			break
		}
		term /= float64(n)
		total += term

		if term < eps {
			break
		}
		n++
	}
	return total, n
}

func PrintTable(rows [][5]string) {
	sep := "+----------------+----------------+----------------+----------------+----------------+"

	fmt.Println(sep)
	fmt.Printf("| %-14s | %-14s | %-14s | %-14s | %-14s |\n",
		"Функция", "Приближённо", "Точно", "Погрешность", "Итераций")
	fmt.Println(sep)

	for _, r := range rows {
		fmt.Printf("| %-14s | %-14s | %-14s | %-14s | %-14s |\n",
			r[0], r[1], r[2], r[3], r[4])
	}
	fmt.Println(sep)
}

func main() {
	eps := 1e-10
	x := 1.2
	maxIter := 1000000

	var rows [][5]string

	approx, it := sinTaylor(x, eps, 10000)
	rows = append(rows, [5]string{
		"sin(1.2)",
		fmt.Sprintf("%.12f", approx),
		fmt.Sprintf("%.12f", math.Sin(x)),
		fmt.Sprintf("%.2e", math.Abs(approx-math.Sin(x))),
		fmt.Sprintf("%d", it),
	})

	approx, it = cosTaylor(x, eps, 10000)
	rows = append(rows, [5]string{
		"cos(1.2)",
		fmt.Sprintf("%.12f", approx),
		fmt.Sprintf("%.12f", math.Cos(x)),
		fmt.Sprintf("%.2e", math.Abs(approx-math.Cos(x))),
		fmt.Sprintf("%d", it),
	})

	approx, it = piNilakantha(eps, maxIter)
	rows = append(rows, [5]string{
		"π (Нилаканта)",
		fmt.Sprintf("%.12f", approx),
		fmt.Sprintf("%.12f", math.Pi),
		fmt.Sprintf("%.2e", math.Abs(approx-math.Pi)),
		fmt.Sprintf("%d", it),
	})

	approx, it = eTaylor(eps, maxIter)
	rows = append(rows, [5]string{
		"e",
		fmt.Sprintf("%.12f", approx),
		fmt.Sprintf("%.12f", math.E),
		fmt.Sprintf("%.2e", math.Abs(approx-math.E)),
		fmt.Sprintf("%d", it),
	})

	PrintTable(rows)
}
