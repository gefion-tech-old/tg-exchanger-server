package models

import (
	"regexp"

	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ AppValidation.ResourceI = (*Bill)(nil)
var _ AppValidation.ResourceI = (*RejectBill)(nil)

type Bill struct {
	ID        int    `json:"id"`
	ChatID    int64  `json:"chat_id"`
	Bill      string `json:"bill"`
	CreatedAt string `json:"created_at"`
}

type RejectBill struct {
	ChatID int64  `json:"chat_id"`
	Bill   string `json:"bill"`
	Reason string `json:"reason"`
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (rb *RejectBill) Validation() error {
	return validation.ValidateStruct(
		rb,
		validation.Field(&rb.ChatID, validation.Required),
		validation.Field(
			&rb.Bill,
			validation.Required,
			validation.Length(16, 16),
			validation.Match(regexp.MustCompile(AppValidation.RegexCard)),
		),
		validation.Field(&rb.Reason, validation.Required),
	)
}

func (b *Bill) Validation() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.ID,
			validation.When(b.CreatedAt != "",
				validation.Required,
				validation.Min(1),
			),
		),

		validation.Field(&b.ChatID, validation.Required),

		validation.Field(&b.Bill,
			validation.Required,
			validation.Length(16, 16),
			validation.Match(regexp.MustCompile(AppValidation.RegexCard)),
		),

		validation.Field(&b.CreatedAt, validation.When(b.ID > 0,
			validation.Required,
			validation.By(AppValidation.DateValidation(b.CreatedAt))),
		),
	)
}
