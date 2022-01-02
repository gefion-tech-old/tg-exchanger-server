package logs

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

func (m *ModLogs) DeleteLogRecordHandler(c *gin.Context) {
	if obj := m.responser.RecordHandler(c, &models.LogRecord{}).(*models.LogRecord); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.LogRecord{}) {
			return
		}

		m.responser.DeleteRecordResponse(c, m.repository, obj)
		return
	}

	m.responser.Error(c, http.StatusInternalServerError, errors.ErrFailedToInitializeStruct)
}

func (m *ModLogs) GetLogRecordsSelectionHandler(c *gin.Context) {
	m.responser.SelectionResponse(c, m.repository, func(arr interface{}) (interface{}, int) {
		if c.Query("user") != "" {
			newArr := []*models.LogRecord{}
			for _, v := range arr.([]*models.LogRecord) {
				if v.Username != nil {
					if *v.Username == c.Query("user") {
						newArr = append(newArr, v)
					}
				}
			}

			arr = newArr
		}

		dt := "2022-01-02"
		dt1, err := time.Parse("2006-01-02", dt)
		if err != nil {
			fmt.Println(err)
		}

		dt2 := "2022-01-01"
		dt3, err := time.Parse("2006-01-02", dt2)
		if err != nil {
			fmt.Println(err)
		}

		dt3.Before(dt1)

		return arr, len(arr.([]*models.LogRecord))

	})
}

func (m *ModLogs) DeleteLogRecordsSelectionHandler(c *gin.Context) {}
