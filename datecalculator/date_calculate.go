package datecalculator

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	parts := strings.Fields(repeat)

	err := checkErr(parts)
	if err != nil {
		return "", err
	}

	originalDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	var resDate string
	switch parts[0] {
	case "d":
		resDate, err = calculateNextDateByDays(now, originalDate, parts)
	case "y":
		resDate, err = calculateNextDateByYears(now, originalDate), nil
	default:
		return "", errors.New("неверный формат правила повторения события")
	}

	return resDate, err
}

func checkErr(parts []string) error {
	if parts[0] == "d" && len(parts) < 2 {
		return errors.New("указаны не все параметры правила повторения события")
	}
	if parts[0] == "d" {
		numForDay, err := strconv.Atoi(parts[1])
		if err != nil {
			return errors.New("дни указаны не цифрой")
		}
		if numForDay > 400 {
			return errors.New("количество дней больше четырёх сотен")
		}
	}
	if parts[0] == "y" && len(parts) > 1 {
		return errors.New("указано избыточное число параметров для года")
	}
	return nil
}

func calculateNextDateByDays(now, originalDate time.Time, parts []string) (string, error) {
	daysToAdd, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	nextDate := originalDate.AddDate(0, 0, daysToAdd)
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(0, 0, daysToAdd)
	}
	return nextDate.Format("20060102"), nil
}

func calculateNextDateByYears(now, originalDate time.Time) string {
	nextDate := originalDate.AddDate(1, 0, 0)
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format("20060102")
}
