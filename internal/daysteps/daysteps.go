package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать количество шагов: %w", err)
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать продолжительность: %w", err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	daySteps, dayDur, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return " "
	}
	if daySteps <= 0 {
		return " "
	}
	distanceInMeters := float64(daySteps) * StepLength
	kmDistance := distanceInMeters / 1000.0
	Calories := spentcalories.WalkingSpentCalories(daySteps, weight, height, dayDur)
	result := fmt.Sprintf(`Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.`, daySteps, kmDistance, Calories)

	return result
}
