package bills

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gin-gonic/gin"
)

/*
	@Method DELETE
	@Path admin/bill/reject
	@Type PRIVATE
	@Documentation

	Отклонить верификацию карты
*/
func (m *ModBills) RejectBillHandler(c *gin.Context) {
	req := &models.RejectBill{}
	if err := c.ShouldBindJSON(req); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if err := req.RejectBillValidation(); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, err)
		return
	}

	payload, err := json.Marshal(map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id": req.ChatID,
		},
		"message": map[string]interface{}{
			"type": "confirmation_successful",
			"text": fmt.Sprintf("🔴 Отклонение 🔴\n\nКарта `%s` отклонена.\nПричина: %s", req.Bill, req.Reason),
		},
		"created_at": time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
	})
	if err != nil {
		fmt.Println(err)
	}

	if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{})
}

/*
	@Method DELETE
	@Path admin/bill
	@Type PRIVATE
	@Documentation

	Создать запись в таблице `bills`

	# TESTED
*/
func (m *ModBills) CreateBillHandler(c *gin.Context) {
	r := &models.Bill{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.BillValidation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Bill{}) {
			return
		}

		m.responser.CreateRecordResponse(c, m.store.AdminPanel().Bills(), obj)
		payload, err := json.Marshal(map[string]interface{}{
			"to": map[string]interface{}{
				"chat_id": r.ChatID,
			},
			"message": map[string]interface{}{
				"type": "confirmation_successful",
				"text": fmt.Sprintf("🟢 Успех! 🟢\n\nКарта `%s` успешно верифицирован.", r.Bill),
			},
			"created_at": time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
		})
		if err != nil {
			fmt.Println(err)
		}

		if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
			fmt.Println(err)
		}

		return
	}
}
