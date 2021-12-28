package models

import (
	_errors "errors"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	validation "github.com/go-ozzo/ozzo-validation"
)

// Структура записи в таблице `users`
type User struct {
	ChatID    int64   `json:"chat_id"`
	Username  string  `json:"username" binding:"required"`
	Hash      *string `json:"hash"`
	Role      int     `json:"role"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type UserFromBotRequest struct {
	ChatID   int64  `json:"chat_id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UserFromAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Testing  bool   `json:"testing"`
}

type UserCodeRequest struct {
	Code uint64 `json:"code" binding:"required"`
}

/*
	==========================================================================================
	ВАЛИДАЦИЯ ДАННЫХ
	==========================================================================================
*/

func (req *UserCodeRequest) UserCodeRequestValidation() error {
	return validation.ValidateStruct(
		req,
		validation.Field(
			&req.Code,
			validation.By(verificationСodeValidation(req.Code)),
		),
	)
}

func (req *UserFromAdminRequest) UserFromAdminRequestValidation(urs config.UsersConfig) error {
	return validation.ValidateStruct(
		req,
		validation.Field(
			&req.Username,
			validation.By(userRightsValidation(req.Username, urs)),
			validation.Length(3, 20),
		),
		validation.Field(
			&req.Password,
			validation.Length(8, 15),
		),
	)
}

// Функция проверки валидности кода
func verificationСodeValidation(code uint64) validation.RuleFunc {
	return func(value interface{}) error {
		if code >= 100000 && code <= 999999 {
			return nil
		}

		return _errors.New("is invalid")
	}
}

// Проверяем, имеет ли данный пользователь права регестрироваться в админке
func userRightsValidation(uname string, urs config.UsersConfig) validation.RuleFunc {
	return func(value interface{}) error {
		uArr := []string{}
		uArr = append(uArr, urs.Managers...)
		uArr = append(uArr, urs.Developers...)
		uArr = append(uArr, urs.Admins...)

		for _, m := range uArr {
			if uname == m {
				return nil
			}
		}

		return errors.ErrNotEnoughRights
	}
}
