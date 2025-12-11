package time

import (
	"fmt"
	"strconv"
	"time"
)

func CreateUnixTimeFromDays(d string) (int64, error) {
	r, err := ParseDays(d)
	if err != nil {
		return 0, err
	}
	// Дни в секундах плюс один час для просмотра после бана в настройках
	return time.Now().Unix() + (r*86400 + 600), nil
}

func ParseDays(s string) (int64, error) {
	days, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("не удалось распарсить days %q в int64 %w", s, err)
	}

	return days, nil
}
