package models

import (
	"fmt"
	"regexp"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

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
				static.NTF__T__VERIFICATION,
				static.NTF__T__EXCHANGE_ERROR,
				static.NTF__T__REQ_SUPPORT,
			),
		),

		validation.Field(&n.Status,
			validation.Required,
			validation.In(
				static.NTF__S__NEW,
				static.NTF__S__IN_PROCESS,
				static.NTF__S__COMPLETED,
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
				validation.By(DateValidation(n.CreatedAt)),
			).Else(validation.Empty),
		),

		validation.Field(&n.UpdatedAt,
			validation.When(n.CreatedAt != "",
				validation.Required,
				validation.By(DateValidation(n.UpdatedAt)),
			).Else(validation.Empty),
		),
	)
}

func NotificationMetaDataValidation(nmt *NotificationMetaData, nType int) validation.RuleFunc {
	return func(value interface{}) error {
		return validation.ValidateStruct(
			nmt,

			validation.Field(&nmt.CardVerification,
				validation.When(nType == static.NTF__T__VERIFICATION,
					validation.By(CardVerificationMetaDataValidation(&nmt.CardVerification)),
				),
			),

			validation.Field(&nmt.ExActionCancel,
				validation.When(nType == static.NTF__T__EXCHANGE_ERROR,
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
				validation.Min(100000),
				validation.Max(999999),
			),

			validation.Field(&cvmt.ImgPath, validation.Required),

			validation.Field(&cvmt.UserCard,
				validation.Required,
				validation.Match(regexp.MustCompile(static.REGEX__CARD)),
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
	fmt.Println(111)
	return func(value interface{}) error {
		return validation.ValidateStruct(
			nud,
			validation.Field(&nud.ChatID,
				validation.Required,
			),

			validation.Field(&nud.Username,
				validation.Required,
				validation.Match(regexp.MustCompile(static.REGEX__NAME)),
			),
		)
	}
}
