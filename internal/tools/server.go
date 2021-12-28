package tools

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gin-gonic/gin"
)

/*
	Определение какую роль должен получить регистрируемый в админке пользователь
*/
func RoleDefine(uname string, urs config.UsersConfig) int {
	for _, v := range urs.Admins {
		if v == uname {
			return static.S__ROLE__ADMIN
		}
	}

	for _, v := range urs.Developers {
		if v == uname {
			return static.S__ROLE__ADMIN
		}
	}

	for _, v := range urs.Managers {
		if v == uname {
			return static.S__ROLE__MANAGER
		}
	}

	return static.S__ROLE__USER
}

func ServErr(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
	c.Abort()
}
