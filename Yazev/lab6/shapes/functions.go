package shapes

import "fmt"

func PrintShapesInfo(shapes []Shape) {
	fmt.Println(" Информация о фигурах ")
	for i, s := range shapes {
		fmt.Printf("%d) %-40s тип=%-20s площадь=%8.2f периметр=%8.2f\n",
			i+1, s.String(), s.Type(), s.Area(), s.Perimeter())
	}
}

func TotalArea(shapes ...Shape) float64 {
	var total float64
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func FilterShapes(shapes []Shape, pred func(Shape) bool) []Shape {
	var result []Shape
	for _, s := range shapes {
		if pred(s) {
			result = append(result, s)
		}
	}
	return result
}
