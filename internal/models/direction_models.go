package models

import (
	"regexp"

	AppInterfaces "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var _ AppInterfaces.ResourceI = (*Direction)(nil)
var _ AppInterfaces.ResourceI = (*DirectionMA)(nil)
var _ AppInterfaces.ResourceI = (*DirectionSelection)(nil)

type Direction struct {
	ID                  int    `json:"id"`
	ExchangeFrom        string `json:"exchange_from"`
	ExchangeTo          string `json:"exchange_to"`
	CourseCorrection    int    `json:"course_correction"`
	AddressVerification bool   `json:"address_verification"`
	Status              bool   `json:"status"`
	CreatedBy           string `json:"created_by"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type DirectionMA struct {
	ID          int    `json:"id"`
	DirectionID int    `json:"direction_id"`
	MaID        int    `json:"ma_id"`
	ServiceType int    `json:"service_type"`
	Status      bool   `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type DirectionSelection struct {
	Page   *int
	Limit  *int
	Status *bool
}

type DirectionMASelection struct {
	Page        *int
	Limit       *int
	DirectionID int
}

func (dmas *DirectionMASelection) Validation() error {
	return validation.ValidateStruct(
		dmas,
		validation.Field(
			&dmas.Page,
			validation.When(
				dmas.Page != nil,
				validation.Required,
				validation.Min(1)),
		),

		validation.Field(
			&dmas.Limit,
			validation.When(
				dmas.Limit != nil,
				validation.Required,
				validation.Min(1),
				validation.Max(30),
			),
		),

		validation.Field(
			&dmas.DirectionID,
			validation.Required,
			validation.Min(1),
		),
	)
}

func (ds *DirectionSelection) Validation() error {
	return validation.ValidateStruct(
		ds,
		validation.Field(
			&ds.Page,
			validation.When(
				ds.Page != nil,
				validation.Required,
				validation.Min(1)),
		),

		validation.Field(
			&ds.Limit,
			validation.When(
				ds.Limit != nil,
				validation.Required,
				validation.Min(1),
				validation.Max(30),
			),
		),

		validation.Field(
			&ds.Status,
			validation.When(
				ds.Status != nil,
				validation.In(true, false),
			),
		),
	)
}

func (dma *DirectionMA) Validation() error {
	return nil
}

func (d *Direction) Validation() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.ID,
			validation.When(
				d.CreatedAt != "" && d.UpdatedAt != "",
				validation.Required,
			),
		),

		validation.Field(
			&d.ExchangeFrom,
			is.UpperCase,
			validation.Required,
		),

		validation.Field(
			&d.ExchangeTo,
			is.UpperCase,
			validation.Required,
		),

		validation.Field(
			&d.CourseCorrection,
			validation.Required,
		),

		validation.Field(
			&d.AddressVerification,
			validation.In(true, false),
		),

		validation.Field(
			&d.Status,
			validation.In(true, false),
		),

		validation.Field(
			&d.CreatedBy,
			validation.When(d.CreatedBy != "",
				validation.Required,
				validation.Match(
					regexp.MustCompile(AppValidation.RegexName),
				),
			),
		),

		validation.Field(&d.CreatedAt,
			validation.When(d.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(d.CreatedAt)),
			).Else(validation.Empty),
		),

		validation.Field(&d.UpdatedAt,
			validation.When(d.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(d.UpdatedAt)),
			).Else(validation.Empty),
		),
	)
}
