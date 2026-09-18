package main

import "fmt"

type employee struct {
	name    []byte
	surName []byte
	id      uint8
	salary  uint32
}

type people struct {
	name           []byte
	surName        []byte
	age            uint8
	job            []byte
	employeeSostav []employee
}

func (p people) talk() {
	fmt.Println("Привет, я", string(p.name), string(p.surName))
	fmt.Println("Мне", p.age, "лет и я", string(p.job))
	fmt.Println()
}

func (p people) info() {
	fmt.Println("Вот информация о составе моей команды:")
	if len(p.employeeSostav) == 0 {
		fmt.Println("  (состав пуст)")
		fmt.Println()
		return
	}
	for i, e := range p.employeeSostav {
		fmt.Printf("  %d) %s %s, ID=%d, зарплата=%d\n",
			i+1, string(e.name), string(e.surName), e.id, e.salary)
	}
	fmt.Println()
}

func (p *people) afterFiveYears() {
	p.age += 5
	p.job = []byte("владею компанией")
}

func (p *people) addEmployee(e employee) {
	p.employeeSostav = append(p.employeeSostav, e)
}

func main() {
	andrey := people{
		name:    []byte("Андрей"),
		surName: []byte("Язев"),
		age:     20,
		job:     []byte("безработный"),
		employeeSostav: []employee{
			{name: []byte("Петя"), surName: []byte("Иванов"), id: 1, salary: 50000},
			{name: []byte("Вася"), surName: []byte("Петров"), id: 2, salary: 60000},
		},
	}
	sereja := people{
		name:    []byte("Сережа"),
		surName: []byte("Сясин"),
		age:     30,
		job:     []byte("безработный"),
		employeeSostav: []employee{
			{name: []byte("Коля"), surName: []byte("Сидоров"), id: 3, salary: 70000},
		},
	}

	fmt.Printf("\n\n\n\nОго, это что? История успеха? Ахах\n\n")

	andrey.talk()
	sereja.talk()

	fmt.Println("Спустя 5 лет")
	andrey.afterFiveYears()
	sereja.afterFiveYears()

	andrey.talk()
	sereja.talk()

	andrey.info()
	sereja.info()

	andrey.addEmployee(employee{name: []byte("Дима"), surName: []byte("Кузнецов"), id: 4, salary: 55000})
	sereja.addEmployee(employee{name: []byte("Оля"), surName: []byte("Смирнова"), id: 5, salary: 80000})

	fmt.Println("После найма новых сотрудников:")
	andrey.info()
	sereja.info()
}
