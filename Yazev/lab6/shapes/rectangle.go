package shapes

import "fmt"

//класс
type Rectangle struct {
	W, H float64
}

//фабрика
func NewRectangle(w, h float64) *Rectangle {
	if w <= 0 || h <= 0 {
		panic("NewRectangle: стороны должны быть > 0")
	}
	return &Rectangle{W: w, H: h}
}

func (r Rectangle) Area() float64 {
	return r.W * r.H
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.W + r.H)
}

func (r Rectangle) Type() string {
	return "Прямоугольничек"
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Прямоугольничек{%.2f %.2f}", r.W, r.H)
}

func (r *Rectangle) Scale(k float64) {
	if k <= 0 {
		panic("Rectangle.Scale: коэффициент должен быть > 0")
	}
	r.W *= k
	r.H *= k
}
