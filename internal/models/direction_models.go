package models

import AppInterfaces "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"

var _ AppInterfaces.ResourceI = (*Direction)(nil)

type Direction struct {
	ID                  int    `json:"id"`
	ExchangeFrom        string `json:"exchange_from"`
	ExchangeTo          string `json:"exchange_to"`
	CourseCorrection    int    `json:"course_correction"`
	AddressVerification bool   `json:"address_verification"`
	DirectionStatus     bool   `json:"direction_status"`
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

func (d *Direction) Validation() error {
	return nil
}
