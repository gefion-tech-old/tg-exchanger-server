package models

import (
	"encoding/xml"
	"regexp"

	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation"
)

var _ AppValidation.ResourceI = (*ExchangerSelection)(nil)

type Exchanger struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	UrlToParse string `json:"url"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type ExchangerSelection struct {
	Page  int
	Limit int
}

type OneObmen struct {
	XMLName xml.Name       `xml:"rates"`
	Rates   []OneObmenItem `xml:"item"`
}

type OneObmenItem struct {
	XMLName   xml.Name `xml:"item"`
	From      string   `xml:"from"`
	To        string   `xml:"to"`
	In        float64  `xml:"in"`
	Out       float64  `xml:"out"`
	Amount    float64  `xml:"amount"`
	MinAmount string   `xml:"minamount"`
	MaxAmount string   `xml:"maxamount"`
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (es *ExchangerSelection) Validation() error {
	return validation.ValidateStruct(
		es,
		validation.Field(&es.Page,
			validation.Required,
			validation.Min(1),
		),

		validation.Field(&es.Limit,
			validation.Required,
			validation.Min(1),
			validation.Max(30),
		),
	)
}

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
			validation.Match(regexp.MustCompile(AppValidation.REGEX__NAME)),
		),
		validation.Field(
			&e.UrlToParse,
			validation.Required,
			validation.Required,
			validation.Length(3, 255),
			validation.Match(regexp.MustCompile(AppValidation.REGEX__URL)),
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
			validation.Match(regexp.MustCompile(AppValidation.REGEX__NAME)),
		),
		validation.Field(
			&e.UrlToParse,
			validation.Required,
			validation.Length(3, 255),
			validation.Match(regexp.MustCompile(AppValidation.REGEX__URL)),
		),
	)
}
