package tools

import (
	"math"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
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

// Создать хеш из пароля
func EncryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Сравнить хеш и пароль
func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Сгенерировать код подтверждения
func VerificationCode(testing bool) int {
	if testing {
		return 100000
	} else {
		return RandInt(100000, 999999)
	}
}

// Сделать срез запрашиваемых ресурсов
// Необходимо для расчета, какие записи
// отдавать для какой страницы
func UpperThreshold(page, limit, count int) int {
	if page*limit <= count {
		return page * limit
	}
	return count
}

func LowerThreshold(page, limit, count int) int {
	if math.Ceil(float64(count)/float64(limit)) <= float64(page) {
		return int(math.Ceil(float64(count) / float64(limit)))
	}
	return page
}

func OffsetThreshold(page, limit int) int {
	if page > 1 {
		return (page - 1) * limit
	}

	return page - 1
}
