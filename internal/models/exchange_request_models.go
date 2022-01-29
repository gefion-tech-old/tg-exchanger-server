package models

import (
	AppInterfaces "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ AppInterfaces.ResourceI = (*ExchangeRequest)(nil)
var _ AppInterfaces.ResourceI = (*ExchangeRequestSelection)(nil)

type ExchangeRequest struct {
	ID                int                           `json:"id"`
	Status            AppType.ExchangeRequestStatus `json:"request_status"`
	ExchangeFrom      string                        `json:"exchange_from"`
	ExchangeTo        string                        `json:"exchange_to"`
	Course            string                        `json:"course"`
	Address           string                        `json:"address"`
	ExpectedAmount    float64                       `json:"expected_amount"`
	TransferredAmount float64                       `json:"transferred_amount"`
	TransactionHash   *string                       `json:"transaction_hash"`
	CreatedBy         UserFromBotRequest            `json:"created_by"`
	CreatedAt         string                        `json:"created_at"`
	UpdatedAt         string                        `json:"updated_at"`
}

type ExchangeRequestSelection struct {
	Status AppType.ExchangeRequestStatus
	Page   *int
	Limit  *int
}

func (ers *ExchangeRequestSelection) Validation() error {
	return validation.ValidateStruct(
		ers,

		validation.Field(
			&ers.Page,
			validation.When(ers.Page != nil,
				validation.Required,
				validation.Min(1)),
		),

		validation.Field(
			&ers.Limit,
			validation.When(ers.Limit != nil,
				validation.Required,
				validation.Min(1),
				validation.Max(30),
			),
		),

		validation.Field(
			&ers.Status,
			validation.Required,
			validation.In(
				AppType.ExchangeRequestNew,
			),
		),
	)
}

func (er *ExchangeRequest) Validation() error {
	return validation.ValidateStruct(
		er,

		validation.Field(
			&er.ID,
			validation.When(
				er.CreatedAt != "" && er.UpdatedAt != "",
				validation.Required,
			),
		),

		validation.Field(
			&er.Status,
			validation.Required,
			validation.In(
				AppType.ExchangeRequestNew,
			),
		),

		validation.Field(
			&er.ExchangeFrom,
			validation.Required,
		),

		validation.Field(
			&er.ExchangeTo,
			validation.Required,
		),

		validation.Field(
			&er.Course,
			validation.Required,
		),

		// validation.Field(
		// 	&er.CreatedBy,
		// 	validation.Required,
		// 	validation.Match(regexp.MustCompile(AppValidation.RegexName)),
		// ),

		validation.Field(&er.CreatedAt,
			validation.When(er.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(er.CreatedAt)),
			).Else(validation.Empty),
		),

		validation.Field(&er.UpdatedAt,
			validation.When(er.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(er.UpdatedAt)),
			).Else(validation.Empty),
		),
	)
}
