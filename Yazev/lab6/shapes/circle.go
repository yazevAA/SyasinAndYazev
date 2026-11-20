package shapes

import (
	"fmt"
	"math"
)

// класс
type Circle struct {
	R float64
}

func NewCircle(r float64) *Circle {
	if r <= 0 {
		panic("NewCircle: радиус должен быть > 0")
	}
	return &Circle{R: r}
}

func (c Circle) Area() float64 {
	return math.Pi * c.R * c.R
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.R
}

func (c Circle) Type() string {
	return "Кружок"
}

func (c Circle) String() string {
	return fmt.Sprintf("Кружок{r=%.2f}", c.R)
}

func (c *Circle) Scale(k float64) {
	if k <= 0 {
		panic("Circle.Scale: коэффициент должен быть > 0")
	}
	c.R *= k
}
