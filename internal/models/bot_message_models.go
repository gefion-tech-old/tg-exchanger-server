package models

import (
	"regexp"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

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
				validation.Match(regexp.MustCompile(static.REGEX__NAME)),
			).Else(validation.Empty),
		),

		validation.Field(&b.MessageText, validation.Required),

		validation.Field(&b.CreatedBy,
			validation.When(b.Connector != "",
				validation.Required,
				validation.Match(regexp.MustCompile(static.REGEX__NAME)),
			).Else(validation.Empty),
		),

		validation.Field(&b.CreatedAt,
			validation.When(b.CreatedAt != "",
				validation.Required,
				validation.By(DateValidation(b.CreatedAt)),
			).Else(validation.Empty),
		),

		validation.Field(&b.UpdatedAt,
			validation.When(b.CreatedAt != "",
				validation.Required,
				validation.By(DateValidation(b.UpdatedAt)),
			).Else(validation.Empty),
		),
	)
}
