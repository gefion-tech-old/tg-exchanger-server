package utils

import (
	"database/sql"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gin-gonic/gin"
)

type Responser struct{}

type ResponserI interface {
	NewRecord(c *gin.Context, data interface{}, err error)
	Record(c *gin.Context, data interface{}, err error)
	Error(c *gin.Context, code int, err ...error)
}

func InitResponser() ResponserI {
	return &Responser{}
}

/*
	Метод создания записи в БД и автоматической обработкой результата
*/
func (u *Responser) NewRecord(c *gin.Context, data interface{}, err error) {
	switch err {
	case nil:
		c.JSON(http.StatusCreated, data)
	case sql.ErrNoRows:
		u.Error(c, http.StatusNotFound, err)
		return
	default:
		u.Error(c, http.StatusInternalServerError, err)
		return
	}
}

/*
	Метод действия с записью в БД и автоматической обработкой результата
*/
func (u *Responser) Record(c *gin.Context, data interface{}, err error) {
	switch err {
	case nil:
		c.JSON(http.StatusOK, data)
	case sql.ErrNoRows:
		u.Error(c, http.StatusNotFound, errors.ErrAlreadyExists)
		return
	default:
		u.Error(c, http.StatusInternalServerError, err)
		return
	}
}

func (u *Responser) Error(c *gin.Context, code int, errs ...error) {
	for _, err := range errs {
		if err != nil {
			c.JSON(code, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		}
	}
}
