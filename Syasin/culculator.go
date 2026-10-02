package main
import (
 "bufio"
 "fmt"
 "math"
 "os"
)

// Car хранит основные характеристики гоночного автомобиля.
type Car struct {
 Mass              float64 // масса автомобиля, кг
 Power             float64 // мощность двигателя, кВт
 MaxSpeed          float64 // максимальная скорость, км/ч
 AccelerationTime  float64 // время разгона 0-100 км/ч, с
 GripCoefficient   float64 // коэффициент сцепления с дорогой
 TurnRadius        float64 // радиус поворота, м
 EngineRPM         float64 // обороты двигателя, об/мин
 GearRatio         float64 // передаточное число
 WheelRadius       float64 // радиус колеса, м
 DragCoefficient   float64 // коэффициент аэродинамического сопротивления
 FrontalArea       float64 // площадь лобовой поверхности, м^2
 AirDensity        float64 // плотность воздуха, кг/м^3
}

// AverageAcceleration рассчитывает среднее ускорение при разгоне от 0 до 100 км/ч.
func (c Car) AverageAcceleration() float64 {
 speed := 100.0 / 3.6
 return speed / c.AccelerationTime
}

// BrakingDistance рассчитывает примерный тормозной путь.
func (c Car) BrakingDistance() float64 {
 speed := c.MaxSpeed / 3.6
 g := 9.81

 return math.Pow(speed, 2) / (2 * c.GripCoefficient * g)
}

// CorneringSpeed рассчитывает максимальную скорость автомобиля в повороте.
func (c Car) CorneringSpeed() float64 {
 g := 9.81

 speed := math.Sqrt( c.GripCoefficient * g * c.TurnRadius, )

 return speed * 3.6
}

// Torque рассчитывает крутящий момент двигателя.
func (c Car) Torque() float64 {
 power := c.Power * 1000

 angularSpeed := 2 * math.Pi * c.EngineRPM / 60

 return power / angularSpeed
}

// WheelSpeed рассчитывает скорость автомобиля по оборотам двигателя.
func (c Car) WheelSpeed() float64 {
 wheelRPM := c.EngineRPM / c.GearRatio

 wheelLength := 2 * math.Pi * c.WheelRadius

 speedMetersPerMinute := wheelRPM * wheelLength

 speedKmPerHour := speedMetersPerMinute * 60 / 1000

 return speedKmPerHour
}

// DragForce рассчитывает силу аэродинамического сопротивления.
func (c Car) DragForce() float64 {
 speed := c.MaxSpeed / 3.6

 return 0.5 *
  c.AirDensity *
  c.DragCoefficient *
  c.FrontalArea *
  math.Pow(speed, 2)
}

// Downforce рассчитывает примерную прижимную силу.
func (c Car) Downforce() float64 {
 speed := c.MaxSpeed / 3.6

 downforceCoefficient := 1.2

 return 0.5 *
  c.AirDensity *
  downforceCoefficient *
  c.FrontalArea *
  math.Pow(speed, 2)
}

// SteeringAngle рассчитывает угол поворота автомобиля.
func (c Car) SteeringAngle() float64 {
 wheelBase := 2.5

 angleRadians := math.Atan(wheelBase / c.TurnRadius)

 angleDegrees := angleRadians * 180 / math.Pi

 return angleDegrees
}

// PowerToWeight показывает отношение мощности к массе.
func (c Car) PowerToWeight() float64 {
 return c.Power / (c.Mass / 1000)
}

// GripForce рассчитывает максимальную силу сцепления.
func (c Car) GripForce() float64 {
 g := 9.81

 return c.GripCoefficient * c.Mass * g
}

// CarRating оценивает автомобиль по нескольким параметрам.
func (c Car) CarRating() string {
 powerToWeight := c.PowerToWeight()

 if powerToWeight < 100 {
  return "Начальный уровень"
 } else if powerToWeight < 180 {
  return "Хороший уровень"
 } else if powerToWeight < 250 {
  return "Высокий уровень"
 }

 return "Очень высокий уровень"
}

// PrintInfo выводит основные характеристики автомобиля.
func (c Car) PrintInfo() {
 fmt.Println("\n========== ХАРАКТЕРИСТИКИ АВТОМОБИЛЯ ==========")

 fmt.Printf("Масса:                 %.1f кг\n", c.Mass)
 fmt.Printf("Мощность:              %.1f кВт\n", c.Power)
 fmt.Printf("Максимальная скорость: %.1f км/ч\n", c.MaxSpeed)
 fmt.Printf("Разгон 0-100:          %.2f с\n", c.AccelerationTime)
 fmt.Printf("Сцепление:             %.2f\n", c.GripCoefficient)
 fmt.Printf("Радиус поворота:       %.1f м\n", c.TurnRadius)
 fmt.Printf("Обороты двигателя:     %.0f об/мин\n", c.EngineRPM)
 fmt.Printf("Передаточное число:    %.2f\n", c.GearRatio)
 fmt.Printf("Радиус колеса:         %.2f м\n", c.WheelRadius)

 fmt.Println("===============================================")
}
// printMenu выводит меню программы.
func printMenu() {
 fmt.Println("\n============== ГОНОЧНЫЙ ИНЖЕНЕР ==============")
 fmt.Println("1. Среднее ускорение")
 fmt.Println("2. Тормозной путь")
 fmt.Println("3. Максимальная скорость в повороте")
 fmt.Println("4. Крутящий момент двигателя")
 fmt.Println("5. Скорость по оборотам двигателя")
 fmt.Println("6. Сила аэродинамического сопротивления")
 fmt.Println("7. Аэродинамическая прижимная сила")
 fmt.Println("8. Угол поворота колёс")
 fmt.Println("9. Оценка автомобиля")
 fmt.Println("10. Максимальная сила сцепления")
 fmt.Println("11. Вывести все расчёты")
 fmt.Println("0. Выход")
 fmt.Println("===============================================")
}

// readPositiveFloat считывает положительное число.
func readPositiveFloat(reader *bufio.Reader, message string) float64 {
 for {
  fmt.Print(message)

  var value float64

  _, err := fmt.Fscan(reader, &value)

  if err != nil {
   fmt.Println("Ошибка! Введите число.")
   reader.ReadString('\n')
   continue
  }

  if value <= 0 {
   fmt.Println("Ошибка! Значение должно быть больше нуля.")
   continue
  }

  return value
 }
}

// readChoice считывает номер пункта меню.
func readChoice(reader *bufio.Reader) int {
 for {
  fmt.Print("Выберите пункт: ")

  var choice int

  _, err := fmt.Fscan(reader, &choice)

  if err != nil {
   fmt.Println("Ошибка! Введите целое число.")
   reader.ReadString('\n')
   continue
  }

  return choice
 }
}

// readCar получает характеристики автомобиля от пользователя.
func readCar(reader *bufio.Reader) Car {
 var car Car

 fmt.Println("Введите характеристики гоночного автомобиля.")
 fmt.Println("Все значения должны быть положительными.\n")

 car.Mass = readPositiveFloat(
  reader,
  "Масса автомобиля (кг): ",
 )

 car.Power = readPositiveFloat(
  reader,
  "Мощность двигателя (кВт): ",
 )

 car.MaxSpeed = readPositiveFloat(
  reader,
  "Максимальная скорость (км/ч): ",
 )

 car.AccelerationTime = readPositiveFloat(
  reader,
  "Время разгона 0-100 (с): ",
 )

 car.GripCoefficient = readPositiveFloat(
  reader,
  "Коэффициент сцепления: ",
 )

 car.TurnRadius = readPositiveFloat(
  reader,
  "Радиус поворота (м): ",
 )

 car.EngineRPM = readPositiveFloat(
  reader,
  "Обороты двигателя (об/мин): ",
 )

 car.GearRatio = readPositiveFloat(
  reader,
  "Передаточное число: ",
 )

 car.WheelRadius = readPositiveFloat(
  reader,
  "Радиус колеса (м): ",
 )

 car.DragCoefficient = readPositiveFloat(
  reader,
  "Коэффициент сопротивления воздуха: ",
 )

 car.FrontalArea = readPositiveFloat(
  reader,
  "Площадь лобовой поверхности (м^2): ",
 )

 car.AirDensity = readPositiveFloat(
  reader,
  "Плотность воздуха (кг/м^3): ",
 )

 return car
}

// printAllCalculations выводит все результаты.
func printAllCalculations(car Car) {
 fmt.Println("\n================ ВСЕ РАСЧЁТЫ ================")

 fmt.Printf(
  "Среднее ускорение:          %.2f м/с^2\n",
  car.AverageAcceleration(),
 )

 fmt.Printf(
  "Тормозной путь:             %.2f м\n",
  car.BrakingDistance(),
 )

 fmt.Printf(
  "Скорость в повороте:        %.2f км/ч\n",
  car.CorneringSpeed(),
 )

 fmt.Printf(
  "Крутящий момент:            %.2f Н*м\n",
  car.Torque(),
 )

 fmt.Printf(
  "Скорость по оборотам:       %.2f км/ч\n",
  car.WheelSpeed(),
 )

 fmt.Printf(
  "Аэродинамическое сопротивление: %.2f Н\n",
  car.DragForce(),
 )

 fmt.Printf(
  "Прижимная сила:             %.2f Н\n",
  car.Downforce(),
 )

 fmt.Printf(
  "Угол поворота колёс:        %.2f градусов\n",
  car.SteeringAngle(),
 )

 fmt.Printf(
  "Мощность на тонну:          %.2f кВт/т\n",
  car.PowerToWeight(),
 )

 fmt.Printf(
  "Максимальная сила сцепления: %.2f Н\n",
  car.GripForce(),
 )

 fmt.Printf(
  "Оценка автомобиля:          %s\n",
  car.CarRating(),
 )

 fmt.Println("=============================================")
}

func main() {
 reader := bufio.NewReader(os.Stdin)
fmt.Println("==============================================")
 fmt.Println("       ГОНОЧНЫЙ ИНЖЕНЕР - Go")
 fmt.Println("==============================================")
 fmt.Println("Программа для расчёта характеристик автомобиля.")
 fmt.Println()

 car := readCar(reader)

 car.PrintInfo()

 for {
  printMenu()

  choice := readChoice(reader)

  switch choice {

  case 1:
   fmt.Printf(
    "\nСреднее ускорение: %.2f м/с^2\n",
    car.AverageAcceleration(),
   )

  case 2:
   fmt.Printf(
    "\nТормозной путь: %.2f м\n",
    car.BrakingDistance(),
   )

  case 3:
   fmt.Printf(
    "\nМаксимальная скорость в повороте: %.2f км/ч\n",
    car.CorneringSpeed(),
   )

  case 4:
   fmt.Printf(
    "\nКрутящий момент двигателя: %.2f Н*м\n",
    car.Torque(),
   )

  case 5:
   fmt.Printf(
    "\nСкорость по оборотам двигателя: %.2f км/ч\n",
    car.WheelSpeed(),
   )

  case 6:
   fmt.Printf(
    "\nСила аэродинамического сопротивления: %.2f Н\n",
    car.DragForce(),
   )

  case 7:
   fmt.Printf(
    "\nАэродинамическая прижимная сила: %.2f Н\n",
    car.Downforce(),
   )

  case 8:
   fmt.Printf(
    "\nУгол поворота колёс: %.2f градусов\n",
    car.SteeringAngle(),
   )

  case 9:
   fmt.Printf(
    "\nОценка автомобиля: %s\n",
    car.CarRating(),
   )

  case 10:
   fmt.Printf(
    "\nМаксимальная сила сцепления: %.2f Н\n",
    car.GripForce(),
   )

  case 11:
   printAllCalculations(car)

  case 0:
   fmt.Println("\nПрограмма завершена.")
   return

  default:
   fmt.Println("\nОшибка! Такого пункта меню нет.")
  }
 }
}

/*
Porshe 911
Масса автомобиля (кг): 1450
Мощность двигателя (кВт): 386
Максимальная скорость (км/ч): 296
Время разгона 0-100 (с): 3.2
Коэффициент сцепления: 1.2
Радиус поворота (м): 50
Обороты двигателя (об/мин): 9000
Передаточное число: 3.0
Радиус колеса (м): 0.36
Коэффициент сопротивления воздуха: 0.30
Площадь лобовой поверхности (м^2): 2.0
Плотность воздуха (кг/м^3): 1.225

ВАЗ
Масса автомобиля (кг): 1060 
Мощность двигателя (кВт): 56 кВт 
Максимальная скорость (км/ч): 155
Время разгона 0-100 (с): 16.0 
Коэффициент сцепления: 1.2
Радиус поворота (м): 50
Обороты двигателя (об/мин): 5600 
Передаточное число: 3.9
Радиус колеса (м): 0.28 
Коэффициент сопротивления воздуха: 0.52 
Площадь лобовой поверхности (м^2): 1.885 
Плотность воздуха (кг/м^3): 1.225
*/