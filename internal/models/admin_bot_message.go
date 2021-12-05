package models

import validation "github.com/go-ozzo/ozzo-validation"

type BotMessage struct {
	Connector string `json:"connector"`
	Text      string `json:"text"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
		validation.Field(&b.Text, validation.Required),
		validation.Field(&b.CreatedBy, validation.Required, validation.By(userRightsValidation(b.CreatedBy, managers, developers))),
	)
}
