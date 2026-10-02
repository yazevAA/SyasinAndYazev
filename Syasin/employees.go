package main

import (
	"fmt"
	"sort"
)

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
}

func (e Employee) ToString() string {
	return fmt.Sprintf("| %-4d | %-20s | %-18s | %10.2f |",
		e.ID,
		e.Name,
		e.Department,
		e.Salary,
	)
}

func addEmployee(employees []Employee, employee Employee) []Employee {
	return append(employees, employee)
}

func filterByDepartment(employees []Employee, department string) []Employee {
	var result []Employee

	for _, employee := range employees {
		if employee.Department == department {
			result = append(result, employee)
		}
	}

	return result
}

func averageSalary(employees []Employee) float64 {
	if len(employees) == 0 {
		return 0
	}

	var sum float64

	for _, employee := range employees {
		sum += employee.Salary
	}

	return sum / float64(len(employees))
}

func filterAboveAverageSalary(employees []Employee) []Employee {
	var result []Employee
	average := averageSalary(employees)

	for _, employee := range employees {
		if employee.Salary > average {
			result = append(result, employee)
		}
	}

	return result
}

func sortBySalary(employees []Employee) {
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].Salary < employees[j].Salary
	})
}

func printEmployees(employees []Employee) {
	fmt.Println("+------+----------------------+--------------------+------------+")
	fmt.Println("| ID   | Имя                  | Отдел              | Зарплата   |")
	fmt.Println("+------+----------------------+--------------------+------------+")

	for _, employee := range employees {
		fmt.Println(employee.ToString())
	}

	fmt.Println("+------+----------------------+--------------------+------------+")
}

func main() {
	employees := []Employee{}

	employee1 := Employee{
		ID:         1,
		Name:       "Иван Петров",
		Department: "IT",
		Salary:     95456.59,
	}

	employee2 := Employee{
		ID:         2,
		Name:       "Анна Смирнова",
		Department: "Бухгалтерия",
		Salary:     75984.98,
	}

	employee3 := Employee{
		ID:         3,
		Name:       "Пётр Иванов",
		Department: "IT",
		Salary:     118761.84,
	}

	employee4 := Employee{
		ID:         4,
		Name:       "Мария Кузнецова",
		Department: "Отдел кадров",
		Salary:     78945.97,
	}

	employee5 := Employee{
		ID:         5,
		Name:       "Алексей Волков",
		Department: "IT",
		Salary:     85843.25,
	}


	employees = addEmployee(employees, employee1)
	employees = addEmployee(employees, employee2)
	employees = addEmployee(employees, employee3)
	employees = addEmployee(employees, employee4)
	employees = addEmployee(employees, employee5)

	fmt.Println("Все сотрудники:")
	printEmployees(employees)

	fmt.Printf("\nСредняя зарплата: %.2f\n", averageSalary(employees))

	fmt.Println("\nСотрудники отдела Отдел кадров:")
	printEmployees(filterByDepartment(employees, "Бухгалтерия"))

	fmt.Println("\nЗарплата выше средней:")
	printEmployees(filterAboveAverageSalary(employees))

	fmt.Println("\nСортировка по зарплате:")
	sortBySalary(employees)
	printEmployees(employees)
}