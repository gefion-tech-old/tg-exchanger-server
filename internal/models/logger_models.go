package models

import (
	"fmt"
	"regexp"
	"strconv"

	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ AppValidation.ResourceI = (*LogRecord)(nil)
var _ AppValidation.ResourceI = (*LogRecordSelection)(nil)

type LogRecord struct {
	ID        int     `json:"id"`
	Username  *string `json:"username"`
	Info      string  `json:"info"`
	Service   int     `json:"service"`
	Module    string  `json:"module"`
	CreatedAt string  `json:"created_at"`
}

type LogRecordSelection struct {
	Page     *int
	Limit    *int
	Service  []string
	Username string
	DateFrom string
	DateTo   string
}

func (ls *LogRecordSelection) Validation() error {
	fmt.Println(ls.Page)

	return validation.ValidateStruct(
		ls,
		validation.Field(&ls.Page,
			validation.When(ls.Page != nil,
				validation.Required,
				validation.Min(1)),
		),

		validation.Field(&ls.Limit,
			validation.When(ls.Limit != nil,
				validation.Required,
				validation.Min(1),
				validation.Max(30),
			),
		),

		validation.Field(&ls.Service, validation.When(len(ls.Service) > 0,
			validation.Each(
				validation.In(
					strconv.Itoa(AppType.LogLevelBot),
					strconv.Itoa(AppType.LogLevelServer),
					strconv.Itoa(AppType.LogLevelAdmin),
				),
			),
		).Else(validation.Nil),
		),

		validation.Field(&ls.Username,
			validation.When(ls.Username != "",
				validation.Match(regexp.MustCompile(AppValidation.RegexName)),
			).Else(validation.Empty),
		),

		validation.Field(&ls.DateFrom, validation.When(ls.DateFrom != "",
			validation.By(AppValidation.DateValidation(ls.DateFrom)),
		).Else(validation.Empty)),

		validation.Field(&ls.DateTo, validation.When(ls.DateTo != "",
			validation.By(AppValidation.DateValidation(ls.DateTo)),
		).Else(validation.Empty)),
	)
}

func (l *LogRecord) Validation() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.ID, validation.When(l.ID > 0, validation.Required)),
		validation.Field(&l.Username,
			validation.When(l.Service == AppType.LogLevelAdmin,
				validation.Required,
				validation.Match(regexp.MustCompile(AppValidation.RegexName)),
			).Else(validation.Nil),
		),

		validation.Field(&l.Info, validation.Required),
		validation.Field(&l.Service, validation.Required,
			validation.In(AppType.LogLevelBot, AppType.LogLevelServer, AppType.LogLevelAdmin),
		),

		validation.Field(&l.Module, validation.Required),
		validation.Field(&l.CreatedAt, validation.When(l.ID > 0, validation.Required)),
	)
}
