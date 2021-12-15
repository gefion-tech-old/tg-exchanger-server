package models

import validation "github.com/go-ozzo/ozzo-validation"

type Exchanger struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	UrlToParse string `json:"url"`
	CreatedBy  string `json:"created_by"`
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
