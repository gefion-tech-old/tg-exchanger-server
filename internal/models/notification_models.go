package models

import (
	"regexp"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppInterfaces "github.com/gefion-tech/tg-exchanger-server/internal/core/interfaces"
	AppTypes "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	AppValidation "github.com/gefion-tech/tg-exchanger-server/internal/core/validation"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ AppInterfaces.ResourceI = (*Notification)(nil)
var _ AppInterfaces.ResourceI = (*NotificationSelection)(nil)

type Notification struct {
	ID        int                  `json:"id"`
	Type      int                  `json:"type" binding:"required"`
	Status    int                  `json:"status"`
	MetaData  NotificationMetaData `json:"meta_data"`
	User      NotificationUserData `json:"user"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

type NotificationMetaData struct {
	CardVerification CardVerificationMetaData `json:"card_verification"`
	ExActionCancel   ExActionCancelMetaData   `json:"ex_action_cancel"`
}

type CardVerificationMetaData struct {
	Code     *int    `json:"code"`
	UserCard *string `json:"user_card"`
	ImgPath  *string `json:"img_path"`
}

type ExActionCancelMetaData struct {
	ExFrom *string `json:"ex_from"`
	ExTo   *string `json:"ex_to"`
}

type NotificationUserData struct {
	ChatID   int64  `json:"chat_id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type NotificationSelection struct {
	Page  int
	Limit int
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (ns *NotificationSelection) Validation() error {
	return validation.ValidateStruct(
		ns,
		validation.Field(&ns.Page,
			validation.Required,
			validation.Min(1),
		),

		validation.Field(&ns.Limit,
			validation.Required,
			validation.Min(1),
			validation.Max(30),
		),
	)
}

func (n *Notification) Validation() error {
	return validation.ValidateStruct(
		n,
		validation.Field(&n.ID,
			validation.When(n.CreatedAt != "" && n.UpdatedAt != "",
				validation.Required),
		),

		validation.Field(&n.Type,
			validation.Required,
			validation.In(
				AppTypes.NotifyTypeVerification,
				AppTypes.NotifyTypeExchangeError,
				AppTypes.NotifyTypeReqSupport,
			),
		),

		validation.Field(&n.Status,
			validation.Required,
			validation.In(
				AppTypes.NotifyStatusNew,
				AppTypes.NotifyStatusInProgress,
				AppTypes.NotifyStatusCompleted,
			),
		),

		validation.Field(&n.MetaData,
			validation.Required,
			validation.When(n.Status == 1,
				validation.By(NotificationMetaDataValidation(&n.MetaData, n.Type)),
			),
		),

		validation.Field(&n.User,
			validation.Required,
			validation.By(NotificationUserDataValidation(&n.User)),
		),

		validation.Field(&n.CreatedAt,
			validation.When(n.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(n.CreatedAt)),
			).Else(validation.Empty),
		),

		validation.Field(&n.UpdatedAt,
			validation.When(n.CreatedAt != "",
				validation.Required,
				validation.By(AppValidation.DateValidation(n.UpdatedAt)),
			).Else(validation.Empty),
		),
	)
}

func NotificationMetaDataValidation(nmt *NotificationMetaData, nType int) validation.RuleFunc {
	return func(value interface{}) error {
		return validation.ValidateStruct(
			nmt,

			validation.Field(&nmt.CardVerification,
				validation.When(nType == AppTypes.NotifyTypeVerification,
					validation.By(CardVerificationMetaDataValidation(&nmt.CardVerification)),
				),
			),

			validation.Field(&nmt.ExActionCancel,
				validation.When(nType == AppTypes.NotifyTypeExchangeError,
					validation.By(ExActionCancelMetaDataValidation(&nmt.ExActionCancel)),
				),
			),
		)
	}
}

func CardVerificationMetaDataValidation(cvmt *CardVerificationMetaData) validation.RuleFunc {
	return func(value interface{}) error {
		return validation.ValidateStruct(
			cvmt,
			validation.Field(&cvmt.Code,
				validation.Required,
				validation.Min(core.VerificationCodeMin),
				validation.Max(core.VerificationCodeMax),
			),

			validation.Field(&cvmt.ImgPath, validation.Required),

			validation.Field(&cvmt.UserCard,
				validation.Required,
				validation.Match(regexp.MustCompile(AppValidation.RegexCard)),
			),
		)
	}
}

func ExActionCancelMetaDataValidation(eacmt *ExActionCancelMetaData) validation.RuleFunc {
	return func(value interface{}) error {
		return validation.ValidateStruct(
			eacmt,
			validation.Field(&eacmt.ExFrom,
				validation.Required,
			),

			validation.Field(&eacmt.ExTo,
				validation.Required,
			),
		)
	}
}

func NotificationUserDataValidation(nud *NotificationUserData) validation.RuleFunc {
	return func(value interface{}) error {
		return validation.ValidateStruct(
			nud,
			validation.Field(&nud.ChatID,
				validation.Required,
			),

			validation.Field(&nud.Username,
				validation.Required,
				validation.Match(regexp.MustCompile(AppValidation.RegexName)),
			),
		)
	}
}
