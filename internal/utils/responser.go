package utils

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Responser struct{}

type ResponserI interface {
	NewRecordResponse(c *gin.Context, data interface{}, err error)
	RecordResponse(c *gin.Context, data interface{}, err error)
	SelectionResponse(c *gin.Context, repository interface{})
	DRRhelper(c *gin.Context, model interface{}) interface{}
	DeleteRecordResponse(c *gin.Context, repository, model interface{})
	Error(c *gin.Context, code int, err ...error)
}

func InitResponser() ResponserI {
	return &Responser{}
}

/*
	Метод создания записи в БД и автоматической обработкой результата
*/
func (u *Responser) NewRecordResponse(c *gin.Context, data interface{}, err error) {
	switch err {
	case nil:
		c.JSON(http.StatusCreated, data)
	case sql.ErrNoRows:
		u.Error(c, http.StatusUnprocessableEntity, errors.ErrAlreadyExists)
		return
	default:
		u.Error(c, http.StatusInternalServerError, err)
		return
	}
}

/*
	Метод действия с записью в БД и автоматической обработкой результата
*/
func (u *Responser) RecordResponse(c *gin.Context, data interface{}, err error) {
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

func (u *Responser) DeleteRecordResponse(c *gin.Context, repository, model interface{}) {
	fn, err := GetReflectMethod(repository, "Delete")
	if err != nil {
		u.Error(c, http.StatusInternalServerError, err)
		return
	}

	// u.RecordResponse(c, model)

	retv := fn.Call([]reflect.Value{
		reflect.ValueOf(model),
	})

	if retv[0].Interface() != nil {
		u.Error(c, http.StatusInternalServerError, retv[0].Interface().(error))
		return
	}

	u.RecordResponse(c, model, nil)

}

func (u *Responser) DRRhelper(c *gin.Context, model interface{}) interface{} {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		u.Error(c, http.StatusUnprocessableEntity, err)
		return nil
	}

	// может быть любым типом
	val := reflect.ValueOf(model)

	// если это указатель, разрешаю его значение
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		u.Error(c, http.StatusInternalServerError,
			fmt.Errorf("failed to process the struct %s", reflect.TypeOf(model).String()))
		return nil
	}

	fID := val.FieldByName("ID")
	if !fID.IsValid() && fID.Kind() != reflect.Int {
		u.Error(c, http.StatusInternalServerError,
			fmt.Errorf("in struct %s, field ID is invalid", reflect.TypeOf(model).String()))
		return nil
	}

	fID.SetInt(int64(id))

	return model
}

/*
	Метод для динамической разбивки данных из БД и
	автоматическим HTTP ответом.

	Для использования данного метода у передеваемого репозитория
	должны быть реализованы методы Count и Selection.
*/
func (u *Responser) SelectionResponse(c *gin.Context, repository interface{}) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		u.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		u.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	errs, _ := errgroup.WithContext(c)

	cArr := make(chan interface{})
	cCount := make(chan *int)

	// Подсчет кол-ва записей в таблице
	errs.Go(func() error {
		defer close(cCount)

		fn, err := GetReflectMethod(repository, "Count")
		if err != nil {
			return err
		}

		retv := fn.Call([]reflect.Value{})
		c := int(retv[0].Int())
		if retv[1].Interface() != nil {
			return retv[1].Interface().(error)
		}

		cCount <- &c
		return nil
	})

	// Достаю из БД запрашиваемые записи
	errs.Go(func() error {
		defer close(cArr)

		fn, err := GetReflectMethod(repository, "Selection")
		if err != nil {
			return err
		}

		retv := fn.Call([]reflect.Value{
			reflect.ValueOf(page),
			reflect.ValueOf(limit),
		})

		if retv[1].Interface() != nil {
			return retv[1].Interface().(error)
		}

		cArr <- retv[0].Interface()
		return nil
	})

	arr := <-cArr
	count := <-cCount

	if arr == nil || count == nil {
		u.Error(c, http.StatusInternalServerError, errs.Wait())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"current_page": page,
		"last_page":    math.Ceil(float64(*count) / float64(limit)),
		"data":         arr,
	})
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
