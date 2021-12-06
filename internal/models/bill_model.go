package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Bill struct {
	ID        uint   `json:"id"`
	ChatID    int64  `json:"chat_id"`
	Bill      string `json:"bill"`
	CreatedAt string `json:"created_at"`
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (b *Bill) BillValidation() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.ChatID, validation.Required),
		validation.Field(&b.Bill, validation.Required,
			validation.Match(
				regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`))),
	)
}
