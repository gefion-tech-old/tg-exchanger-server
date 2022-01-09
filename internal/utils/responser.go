package utils

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Responser struct {
	logger   LoggerI
	template models.LogRecord
}

type ResponserI interface {
	NewRecordResponse(c *gin.Context, data interface{}, err error)
	RecordResponse(c *gin.Context, data interface{}, err error)
	SelectionResponse(c *gin.Context, repository, querys interface{}) error
	RecordHandler(c *gin.Context, model interface{}, validators ...error) interface{}
	DeleteRecordResponse(c *gin.Context, repository, model interface{}, todo ...func() error) error
	UpdateRecordResponse(c *gin.Context, repository, model interface{}, todo ...func() error) error
	CreateRecordResponse(c *gin.Context, repository, model interface{}, todo ...func() error) error
	Error(c *gin.Context, code int, err ...error) error
}

func InitResponser(l LoggerI) ResponserI {
	return &Responser{
		logger: l,
		template: models.LogRecord{
			Module:  "HTTP RESPONSER",
			Service: static.L__SERVER,
		},
	}
}

/*
	Метод создания записи в БД и автоматической обработкой результата
*/
func (u *Responser) NewRecordResponse(c *gin.Context, data interface{}, err error) {
	switch err {
	case nil:
		c.JSON(http.StatusCreated, data)
	case sql.ErrNoRows:
		u.Error(c, http.StatusUnprocessableEntity, AppError.ErrAlreadyExists)
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
		u.Error(c, http.StatusNotFound, AppError.ErrRecordNotFound)
		return
	default:
		u.Error(c, http.StatusInternalServerError, err)
		return
	}
}

/*
	Метод для удаления записи любого ресурса, при условии
	что репозиторий ресурса содержит метод Delete.

	Автоматический HTTP ответ.
*/
func (u *Responser) DeleteRecordResponse(c *gin.Context, repository, model interface{}, todo ...func() error) error {
	// Инициализация метода операции с БД
	fn, err := GetReflectMethod(repository, "Delete")
	if err != nil {
		return u.Error(c, http.StatusInternalServerError, err)
	}

	// Выполнение операции с БД
	if obj, err := u.callReflectMethod(c, fn, model); obj != nil {
		u.RecordResponse(c, model, nil)
		return nil
	} else if err != nil {
		u.RecordResponse(c, nil, err)
		return err
	}

	return u.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	Метод для обновления записи любого ресурса, при условии
	что репозиторий ресурса содержит метод Update.

	Автоматический HTTP ответ.
*/
func (u *Responser) UpdateRecordResponse(c *gin.Context, repository, model interface{}, todo ...func() error) error {
	// Инициализация метода операции с БД
	fn, err := GetReflectMethod(repository, "Update")
	if err != nil {
		return u.Error(c, http.StatusInternalServerError, err)
	}

	// Выполнение операции с БД
	if obj, err := u.callReflectMethod(c, fn, model); obj != nil {
		u.RecordResponse(c, model, nil)
		return nil
	} else if err != nil {
		u.RecordResponse(c, nil, err)
		return err
	}

	return u.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	Метод для создания записи любого ресурса, при условии
	что репозиторий ресурса содержит метод Create.

	Автоматический HTTP ответ.
*/
func (u *Responser) CreateRecordResponse(c *gin.Context, repository, model interface{}, todo ...func() error) error {
	// Выполнение всех вложенных методов
	for _, executor := range todo {
		if err := executor(); err != nil {
			return u.Error(c, http.StatusInternalServerError, err)
		}
	}

	// Инициализация метода операции с БД
	fn, err := GetReflectMethod(repository, "Create")
	if err != nil {
		return u.Error(c, http.StatusInternalServerError, err)
	}

	// Выполнение операции с БД
	if obj, err := u.callReflectMethod(c, fn, model); obj != nil {
		u.NewRecordResponse(c, model, nil)
		return nil
	} else if err != nil {
		u.NewRecordResponse(c, nil, err)
		return err
	}

	return u.Error(c, http.StatusInternalServerError, AppError.ErrFailedToInitializeStruct)
}

/*
	Универсальный метод для подготовки и валидирования любой структуры.
*/
func (u *Responser) RecordHandler(c *gin.Context, model interface{}, validators ...error) interface{} {
	if c.Param("id") != "" {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return u.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidPathParams)
		}

		// Может быть любым типом
		val := reflect.ValueOf(model)

		// Если это указатель
		if val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}

		if val.Kind() != reflect.Struct {
			err := fmt.Errorf("failed to process the struct %s", reflect.TypeOf(model).String())
			return u.Error(c, http.StatusInternalServerError, err)
		}

		fID := val.FieldByName("ID")
		if !fID.IsValid() && fID.Kind() != reflect.Int {
			err := fmt.Errorf("in struct %s, field ID is invalid", reflect.TypeOf(model).String())
			return u.Error(c, http.StatusInternalServerError, err)
		}

		fID.SetInt(int64(id))
	}

	if len(validators) > 0 {
		if err := u.Error(c, http.StatusUnprocessableEntity, validators...); err != nil {
			return err
		}
	}
	return model
}

/*
	Метод для динамической разбивки данных из БД и
	автоматическим HTTP ответом.

	Для использования данного метода у передеваемого репозитория
	должны быть реализованы методы Count и Selection.
*/
func (u *Responser) SelectionResponse(c *gin.Context, repository, querys interface{}) error {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		return u.Error(c, http.StatusUnprocessableEntity, err)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		return u.Error(c, http.StatusUnprocessableEntity, err)
	}

	// Заполнение объекта querys
	// Может быть любым типом
	val := reflect.ValueOf(querys)

	// Если это указатель
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		err := fmt.Errorf("failed to process the struct %s", reflect.TypeOf(querys).String())
		return u.Error(c, http.StatusInternalServerError, err)
	}

	fPage := val.FieldByName("Page")
	if !fPage.IsValid() && fPage.Kind() != reflect.Int {
		err := fmt.Errorf("in struct %s, field ID is invalid", reflect.TypeOf(querys).String())
		return u.Error(c, http.StatusInternalServerError, err)
	}

	if fPage.Kind() == reflect.Ptr {
		fPage.Set(reflect.ValueOf(&page))
	} else {
		fPage.SetInt(int64(page))
	}

	fLimit := val.FieldByName("Limit")
	if !fLimit.IsValid() && fLimit.Kind() != reflect.Int {
		err := fmt.Errorf("in struct %s, field ID is invalid", reflect.TypeOf(querys).String())
		return u.Error(c, http.StatusInternalServerError, err)
	}

	if fLimit.Kind() == reflect.Ptr {
		fLimit.Set(reflect.ValueOf(&limit))
	} else {
		fLimit.SetInt(int64(limit))
	}

	// Вызов метода валидации структуры запросов

	vfn, err := GetReflectMethod(querys, "Validation")
	if err != nil {
		return err
	}

	validation := vfn.Call([]reflect.Value{})
	if validation[0].Interface() != nil {
		return u.Error(c, http.StatusUnprocessableEntity, validation[0].Interface().(error))
	}

	// Выполнение операций с БД

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

		retv := fn.Call([]reflect.Value{
			reflect.ValueOf(querys),
		})
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
			reflect.ValueOf(querys),
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
		return u.Error(c, http.StatusInternalServerError, errs.Wait())
	}

	c.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"current_page": page,
		"last_page":    math.Ceil(float64(*count) / float64(limit)),
		"total":        count,
		"data":         arr,
	})
	return nil
}

/*
	Универсальный метод обработки операций которые могут вернуть ошибки.
	Если в результате выполнения операции получена ошибка, запрос автоматически прерывается.

	Если ошибка произошла на стороне сервера, она записывается в логи.
*/
func (u *Responser) Error(c *gin.Context, code int, errs ...error) error {
	for _, err := range errs {
		if err != nil {
			// Если 500 ошибка, записываю ее в логи
			if code == http.StatusInternalServerError {
				record := u.template
				record.Info = err.Error()
				u.logger.NewRecord(&record)
			}
			c.JSON(code, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return err
		}
	}
	return nil
}

/*
	==========================================================================================
	ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
	==========================================================================================
*/

func (u *Responser) callReflectMethod(c *gin.Context, fn *reflect.Value, model interface{}) (interface{}, error) {
	retv := fn.Call([]reflect.Value{
		reflect.ValueOf(model),
	})

	if retv[0].Interface() != nil {
		return nil, retv[0].Interface().(error)
	}

	return model, nil
}
