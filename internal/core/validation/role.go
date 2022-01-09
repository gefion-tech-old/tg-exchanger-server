package cvalidation

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
)

// Определение какую роль должен получить регистрируемый в админке пользователь
func RoleDefine(uname string, urs config.UsersConfig) int {
	for _, v := range urs.Admins {
		if v == uname {
			return AppType.AppRoleAdmin
		}
	}

	for _, v := range urs.Developers {
		if v == uname {
			return AppType.AppRoleAdmin
		}
	}

	for _, v := range urs.Managers {
		if v == uname {
			return AppType.AppRoleManager
		}
	}

	return AppType.AppRoleUser
}
