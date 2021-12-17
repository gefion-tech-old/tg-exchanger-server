package models

import (
	"regexp"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Exchanger struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	UrlToParse string `json:"url"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (e *Exchanger) ExchangerCreateValidation() error {
	return validation.ValidateStruct(
		e,
		validation.Field(
			&e.CreatedBy,
			validation.Required,
			validation.Length(3, 20),
		),
		validation.Field(
			&e.Name,
			validation.Required,
			validation.Length(3, 10),
			validation.Match(regexp.MustCompile(static.REGEX__NAME)),
		),
		validation.Field(
			&e.UrlToParse,
			validation.Required,
			validation.Required,
			validation.Length(3, 255),
			validation.Match(regexp.MustCompile(static.REGEX__URL)),
		),
	)
}

func (e *Exchanger) ExchangerUpdateValidation() error {
	return validation.ValidateStruct(
		e,
		validation.Field(
			&e.Name,
			validation.Required,
			validation.Length(3, 10),
			validation.Match(regexp.MustCompile(static.REGEX__NAME)),
		),
		validation.Field(
			&e.UrlToParse,
			validation.Required,
			validation.Length(3, 255),
			validation.Match(regexp.MustCompile(static.REGEX__URL)),
		),
	)
}
