package models

import (
	"regexp"

	AppInterfaces "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ AppInterfaces.ResourceI = (*MerchantAutopayout)(nil)
var _ AppInterfaces.ResourceI = (*MerchantAutopayoutSelection)(nil)

type MerchantAutopayout struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Service     string `json:"service"`
	ServiceType int    `json:"service_type"`
	Options     string `json:"options"`
	Status      bool   `json:"status"`
	MessageID   int    `json:"message_id"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type MerchantAutopayoutSelection struct {
	Page    *int
	Limit   *int
	Service []string
}

/* Whitebit */

type WhitebitOptionParams struct {
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
	BaseURL   string `json:"base_url"`
}

type WhitebitApiHelper struct {
	PublicKey string
	SecretKey string
	BaseURL   string
}

type WhitebitGetHistory struct {
	TransactionMethod int           `json:"transactionMethod"`
	Ticker            string        `json:"ticker"`
	Address           string        `json:"address"`
	UniqueId          string        `json:"uniqueId"`
	Limit             int           `json:"limit"`
	Offset            int           `json:"offset"`
	Status            []interface{} `json:"status"`
}

/* Mine */

type MineOptionParams struct{}

func (ma *MerchantAutopayoutSelection) Validation() error {
	return validation.ValidateStruct(
		ma,
		validation.Field(&ma.Page,
			validation.When(ma.Page != nil,
				validation.Required,
				validation.Min(1)),
		),

		validation.Field(&ma.Limit,
			validation.When(ma.Limit != nil,
				validation.Required,
				validation.Min(1),
				validation.Max(30),
			),
		),

		validation.Field(
			&ma.Service,
			validation.When(len(ma.Service) > 0,
				validation.Each(
					validation.In(
						AppType.MerchantAutoPayoutMine,
						AppType.MerchantAutoPayoutWhitebit,
					),
				),
			),
		),
	)
}

func (ma *MerchantAutopayout) Validation() error {
	return validation.ValidateStruct(
		ma,
		validation.Field(&ma.ID,
			validation.When(ma.CreatedAt != "" && ma.UpdatedAt != "",
				validation.Required,
			),
		),

		validation.Field(&ma.Name,
			validation.Required,
			validation.Match(regexp.MustCompile(AppValidation.RegexName)),
		),

		validation.Field(&ma.Service,
			validation.When(
				ma.CreatedBy != "",
				validation.Required,
				validation.In(
					AppType.MerchantAutoPayoutMine,
					AppType.MerchantAutoPayoutWhitebit,
				),
				validation.Match(regexp.MustCompile(AppValidation.RegexName)),
			),
		),

		validation.Field(&ma.ServiceType,
			validation.Required,
			validation.In(
				AppType.UseAsMerchant,
				AppType.UseAsAutoPayout,
			),
		),

		validation.Field(&ma.Options,
			validation.Required,
		),

		validation.Field(&ma.Status,
			validation.In(true, false),
		),

		validation.Field(&ma.MessageID, validation.Required),

		validation.Field(&ma.CreatedBy,
			validation.When(ma.Service != "",
				validation.Required,
				validation.Match(regexp.MustCompile(AppValidation.RegexName)),
			),
		),

		validation.Field(&ma.CreatedAt,
			validation.When(ma.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(ma.CreatedAt)),
			).Else(validation.Empty),
		),

		validation.Field(&ma.UpdatedAt,
			validation.When(ma.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(ma.UpdatedAt)),
			).Else(validation.Empty),
		),
	)
}
