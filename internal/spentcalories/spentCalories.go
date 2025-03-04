package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, " ", 0, fmt.Errorf("error in the data format")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, " ", 0, fmt.Errorf("steps cannot be converted")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, " ", 0, fmt.Errorf("can't convert time")
	}

	return steps, parts[1], duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	distanceMeters := float64(steps) * lenStep / mInKm
	return distanceMeters
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps)
	speed := dist / duration.Hours()
	return speed
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, actInfo, duration, err := parseTraining(data)
	if err != nil {
		return err.Error()
	}
	if steps <= 0 || duration.Hours() <= 0 {
		return "incorrect training data"
	}

	dist := distance(steps)
	speed := meanSpeed(steps, duration)

	calories := 0.0
	switch actInfo {
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
	default:
		return "unknown type of training"
	}

	result := fmt.Sprintf(`Тип тренировки: %s
   Длительность: %.2f ч.
   Дистанция: %.2f км.
   Скорость: %.2f км/ч
   Сожгли калорий: %.2f`, actInfo, duration.Hours(), dist, speed, calories)
	return result
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	calories := ((runningCaloriesMeanSpeedMultiplier * speed) - runningCaloriesMeanSpeedShift) * weight
	return calories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	dur := duration.Hours()
	calories := ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * dur * minInH
	return calories

}
