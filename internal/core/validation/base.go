package cvalidation

import (
	"regexp"

	CoreErrors "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Функция для валидации стандартизированного для
// приложения данного формата даты.
func DateValidation(d string) validation.RuleFunc {
	return func(value interface{}) error {
		r, err := regexp.Compile(REGEX__DATE)
		if err != nil {
			return err
		}

		if !r.MatchString(d) {
			return CoreErrors.ErrValidationIndalidDateFormat
		}

		return nil
	}
}
