package models

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	validation "github.com/go-ozzo/ozzo-validation"
)

type LogRecord struct {
	ID        int     `json:"id"`
	Username  *string `json:"username"`
	Info      string  `json:"info"`
	Service   int     `json:"service"`
	Module    string  `json:"module"`
	CreatedAt string  `json:"created_at"`
}

type LogRecordSelection struct {
	Page     int
	Limit    int
	Service  []string
	Username string
	DateFrom string
	DateTo   string
}

func (l *LogRecord) InternalRecordValidation() error {
	return validation.ValidateStruct(
		l,
		validation.Field(
			&l.Info,
			validation.Required,
		),
		validation.Field(
			&l.Service,
			validation.Required,
			validation.In(static.L__BOT, static.L__SERVER),
		),
		validation.Field(
			&l.Module,
			validation.Required,
		),
	)
}

func (l *LogRecord) AdminRecordValidation() error {
	return validation.ValidateStruct(
		l,
		validation.Field(
			&l.Info,
			validation.Required,
		),
		validation.Field(
			&l.Username,
			validation.Required,
		),
		validation.Field(
			&l.Service,
			validation.Required,
			validation.In(static.L__ADMIN),
		),
		validation.Field(
			&l.Module,
			validation.Required,
		),
	)
}
