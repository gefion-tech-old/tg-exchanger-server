package models

import validation "github.com/go-ozzo/ozzo-validation"

type BotMessage struct {
	ID          uint   `json:"id"`
	Connector   string `json:"connector"`
	MessageText string `json:"message_text"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (b *BotMessage) BotMessageValidation(managers, developers []string) error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.Connector, validation.Required),
		validation.Field(&b.MessageText, validation.Required),
		validation.Field(&b.CreatedBy, validation.Required, validation.By(userRightsValidation(b.CreatedBy, managers, developers))),
	)
}
