package models

import (
	_errors "errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (ls *LogRecordSelection) Validation() error {
	fmt.Println(ls.DateTo == "")
	return validation.ValidateStruct(
		ls,
		validation.Field(&ls.Service, validation.When(len(ls.Service) > 0,
			validation.Each(
				validation.In(
					strconv.Itoa(static.L__BOT),
					strconv.Itoa(static.L__SERVER),
					strconv.Itoa(static.L__ADMIN),
				),
			),
		).Else(validation.Nil),
		),

		validation.Field(&ls.Username,
			validation.When(ls.Username != "",
				validation.Match(regexp.MustCompile(static.REGEX__NAME)),
			).Else(validation.Empty),
		),

		validation.Field(&ls.DateFrom, validation.When(ls.DateFrom != "",
			validation.By(DateValidation(ls.DateFrom)),
		).Else(validation.Empty)),

		validation.Field(&ls.DateTo, validation.When(ls.DateTo != "",
			validation.By(DateValidation(ls.DateTo)),
		).Else(validation.Empty)),
	)
}

func (l *LogRecord) Validation() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.ID, validation.When(l.ID > 0, validation.Required)),
		validation.Field(&l.Username,
			validation.When(l.Service == static.L__ADMIN,
				validation.Required,
				validation.Match(regexp.MustCompile(static.REGEX__NAME)),
			).Else(validation.Nil),
		),

		validation.Field(&l.Info, validation.Required),
		validation.Field(&l.Service, validation.Required,
			validation.In(static.L__BOT, static.L__SERVER, static.L__ADMIN),
		),

		validation.Field(&l.Module, validation.Required),
		validation.Field(&l.CreatedAt, validation.When(l.ID > 0, validation.Required)),
	)
}

func DateValidation(d string) validation.RuleFunc {
	return func(value interface{}) error {
		r, err := regexp.Compile(static.REGEX__DATE)
		if err != nil {
			return err
		}

		if !r.MatchString(d) {
			return _errors.New("invalid date format")
		}

		return nil
	}
}
