package models

import (
	"regexp"

	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ AppValidation.ResourceI = (*BotMessage)(nil)
var _ AppValidation.ResourceI = (*BotMessageSelection)(nil)

type BotMessage struct {
	ID          int    `json:"id"`
	Connector   string `json:"connector"`
	MessageText string `json:"message_text"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type BotMessageSelection struct {
	Page  int
	Limit int
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (bms *BotMessageSelection) Validation() error {
	return validation.ValidateStruct(
		bms,
		validation.Field(&bms.Page,
			validation.Required,
			validation.Min(1),
		),

		validation.Field(&bms.Limit,
			validation.Required,
			validation.Min(1),
			validation.Max(30),
		),
	)
}

func (b *BotMessage) Validation() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.ID,
			validation.When(b.CreatedAt != "" && b.UpdatedAt != "",
				validation.Required),
		),

		validation.Field(&b.Connector,
			validation.When(b.CreatedBy != "",
				validation.Required,
				validation.Match(regexp.MustCompile(AppValidation.RegexName)),
			).Else(validation.Empty),
		),

		validation.Field(&b.MessageText, validation.Required),

		validation.Field(&b.CreatedBy,
			validation.When(b.Connector != "",
				validation.Required,
				validation.Match(regexp.MustCompile(AppValidation.RegexName)),
			).Else(validation.Empty),
		),

		validation.Field(&b.CreatedAt,
			validation.When(b.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(b.CreatedAt)),
			).Else(validation.Empty),
		),

		validation.Field(&b.UpdatedAt,
			validation.When(b.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(b.UpdatedAt)),
			).Else(validation.Empty),
		),
	)
}
