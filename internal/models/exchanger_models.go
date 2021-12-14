package models

import validation "github.com/go-ozzo/ozzo-validation"

type Exchanger struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	UrlToParse string `json:"url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (e *Exchanger) ExchangerValidation() error {
	return validation.ValidateStruct(
		e,
		validation.Field(&e.Name, validation.Required),
		validation.Field(&e.UrlToParse, validation.Required),
	)
}

func (e *Exchanger) ExchangerValidationFull() error {
	return validation.ValidateStruct(
		e,
		validation.Field(&e.ID, validation.Required, validation.Min(1)),
		validation.Field(&e.Name, validation.Required),
		validation.Field(&e.UrlToParse, validation.Required),
	)
}
