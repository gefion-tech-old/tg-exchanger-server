package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
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

func (b *BotMessage) StructFullness() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.ID, validation.Required),
		validation.Field(&b.Connector, validation.Required),
		validation.Field(&b.MessageText, validation.Required),
		validation.Field(&b.CreatedBy, validation.Required),
		validation.Field(&b.CreatedAt, validation.Required),
		validation.Field(&b.UpdatedAt, validation.Required),
	)
}

func (b *BotMessage) CreateBotMessageValidation(managers, developers []string) error {
	return validation.ValidateStruct(
		b,
		validation.Field(
			&b.Connector,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[^._ ](?:[\w-]|\.[\w-])+[^._ ]$`)),
		),
		validation.Field(&b.MessageText, validation.Required),
		validation.Field(&b.CreatedBy, validation.Required),
	)
}

func (b *BotMessage) UpdateBotMessageValidation(managers, developers []string) error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.MessageText, validation.Required),
	)
}
