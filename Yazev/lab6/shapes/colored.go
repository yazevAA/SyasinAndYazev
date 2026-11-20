package shapes

import "fmt"

type ColoredShape struct {
	Shape //интерфейс
	Color string
}

func NewColoredShape(s Shape, color string) *ColoredShape {
	if s == nil {
		panic("NewColoredShape: Shape не может быть nil")
	}
	if color == "" {
		panic("NewColoredShape: цвет не может быть пустым")
	}
	return &ColoredShape{Shape: s, Color: color}
}

func (cs ColoredShape) String() string {
	return fmt.Sprintf("%s [Цвет=%s]", cs.Shape.String(), cs.Color)
}
