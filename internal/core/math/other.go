package cmath

import (
	"math/rand"
	"time"
)

// Сгенерировать случайное число
func RandInt(min int, max int) int {
	rand.Seed(time.Now().Unix())
	if min > max {
		return min
	} else {
		return rand.Intn(max-min) + min
	}
}

// Сгенерировать код подтверждения
func VerificationCode(testing bool) int {
	if testing {
		return 100000
	} else {
		return RandInt(100000, 999999)
	}
}

// Определение порога запрашиваемых данных
func OffsetThreshold(page, limit int) int {
	if page > 1 {
		return (page - 1) * limit
	}

	return page - 1
}
