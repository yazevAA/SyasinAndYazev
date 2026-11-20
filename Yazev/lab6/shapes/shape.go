package shapes

type Shape interface {
	Area() float64
	Perimeter() float64
	String() string
	Type() string
}
