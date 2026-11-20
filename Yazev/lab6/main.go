package main

import (
	"fmt"

	"lab5/shapes"
)

func main() {

	circle := shapes.NewCircle(2)
	rect := shapes.NewRectangle(3, 4)
	tri := shapes.NewTriangle(3, 4, 5)
	redCircle := shapes.NewColoredShape(shapes.NewCircle(1.5), "красный")
	blueRect := shapes.NewColoredShape(shapes.NewRectangle(2, 5), "прозрачный")

	all := []shapes.Shape{circle, rect, tri, redCircle, blueRect}

	shapes.PrintShapesInfo(all)

	fmt.Printf("\nСуммарная площадь всех фигур: %.2f\n", shapes.TotalArea(all...))

	big := shapes.FilterShapes(all, func(s shapes.Shape) bool {
		return s.Area() > 10
	})
	fmt.Println("\n Фигуры с площадью > 10 ")
	shapes.PrintShapesInfo(big)

	circle.Scale(2)
	rect.Scale(0.5)
	fmt.Printf("\nПосле масштабирования:\n  %s\n  %s\n", circle, rect)

	fmt.Printf("\nТип Кружочек красный: %s, площадь: %.2f\n", redCircle.Type(), redCircle.Area())
}
