package main

import (
	"fmt"
)

var (
	id      string = "ID"
	name    string = "Имя"
	sallary string = "З/П"
)

var idList = []int{67, 69, 52, 911, 1488, 812, 1312, 417}
var nameList = []string{
	"Пупкин Алеша",
	"Кукушкин Таракан",
	"Дураков Фотий",
	"Несмеянов Весельчак",
	"Криворуков Аполлон",
	"Сясин Сергей",
	"Терушкин Борис",
	"Шавалиев Альберт"}

var sallaryList = []float32{23570.91, 6767.6767, 2004.1512, 150000, 409.1308, 25000, 1634234, 0.01}

func main() {

	fmt.Printf("-----------------------------------------------------\n")
	fmt.Printf("|%4s | %-30s | %-10s |\n", id, name, sallary)
	fmt.Printf("-----------------------------------------------------\n")

	for i := 0; i < len(idList); i++ {
		fmt.Printf("|%4d | %-30s | %10.2f |\n", idList[i], nameList[i], sallaryList[i])
	}

	fmt.Printf("-----------------------------------------------------\n")

}
