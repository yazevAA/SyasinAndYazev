package shapes

import (
	"fmt"
	"math"
)

// класс
type Triangle struct {
	A, B, C float64
}

// фабричная функция (создает новый объект)
func NewTriangle(a, b, c float64) *Triangle {
	if a <= 0 || b <= 0 || c <= 0 {
		Ф
		panic("NewTriangle: стороны должны быть > 0")
	}
	if a+b <= c || a+c <= b || b+c <= a {
		panic("NewTriangle: нарушено неравенство треугольника")
	}
	return &Triangle{A: a, B: b, C: c}
}

// методы
func (t Triangle) Area() float64 {
	p := t.Perimeter() / 2
	return math.Sqrt(p * (p - t.A) * (p - t.B) * (p - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

func (t Triangle) Type() string {
	return "Треугольничек"
}

func (t Triangle) String() string {
	return fmt.Sprintf("Тругольничек{%.2f %.2f %.2f}", t.A, t.B, t.C)
}
