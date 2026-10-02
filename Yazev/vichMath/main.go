package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func openImage(filename string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filename)
	case "darwin": // macOS
		cmd = exec.Command("open", filename)
	case "linux":
		cmd = exec.Command("xdg-open", filename)
	default:
		fmt.Println("ОС не поддерживается, откройте файл вручную")
		return
	}

	_ = cmd.Run()
}

func myFunc(x float64) float64 {
	f := x

	return f
}

func main() {
	// Создание графика
	p := plot.New()

	//Оси и название
	p.Title.Text = "График функции"
	p.X.Label.Text = ""
	p.Y.Label.Text = ""
	p.X.Max = 10
	p.X.Min = -10
	p.Y.Max = 10
	p.Y.Min = -10

	line := plotter.NewFunction(myFunc)

	//Размеры и цвета

	line.Samples = 1000 //количество точек
	line.XMin = -10
	line.XMax = 10
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = vg.Points(1)

	p.Add(line)

	if err := p.Save(13*vg.Inch, 9*vg.Inch, "myFunc.png"); err != nil {
		panic(err)
	}

	openImage("myFunc.png")
}
