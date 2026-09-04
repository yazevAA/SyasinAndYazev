package main

import (
    "fmt"
)

// Employee структура для хранения данных о сотруднике
type Employee struct {
    ID         int
    Name       string
    Department string
    Salary     float64
}

// ToString форматированный вывод сотрудника в строку таблицы
func (e Employee) ToString() string {
    return fmt.Sprintf("| %-4d | %-20s | %-15s | %10.2f |", e.ID, e.Name, e.Department, e.Salary)
}

var nextID = 1

// AddEmployee добавляет нового сотрудника в срез
func AddEmployee(employees *[]Employee, name, dept string, salary float64) {
    *employees = append(*employees, Employee{nextID, name, dept, salary})
    nextID++
}

// PrintEmployees выводит всех сотрудников в виде таблицы
func PrintEmployees(employees []Employee, title string) {
    fmt.Println("\n" + title)
    fmt.Println("+------+----------------------+-----------------+------------+")
    fmt.Println("| ID   | Name                 | Department      | Salary     |")
    fmt.Println("+------+----------------------+-----------------+------------+")
    for _, e := range employees {
        fmt.Println(e.ToString())
    }
    fmt.Println("+------+----------------------+-----------------+------------+")
}

func main() {
    employees := []Employee{}

    // Добавление тестовых данных
    AddEmployee(&employees, "Иванов Иван", "IT", 75000)
    AddEmployee(&employees, "Петрова Мария", "HR", 68000)
    AddEmployee(&employees, "Сидоров Сергей", "IT", 82000)
    AddEmployee(&employees, "Козлова Анна", "Finance", 90000)
    AddEmployee(&employees, "Кузнецов Алексей", "HR", 60000)

    // Вывод всех сотрудников
    PrintEmployees(employees, "Список всех сотрудников:")

}