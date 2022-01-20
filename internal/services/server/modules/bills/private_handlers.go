package bills

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/core"
	AppError "github.com/gefion-tech/tg-exchanger-server/internal/core/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gin-gonic/gin"
)

/*
	@Method DELETE
	@Path admin/bill/reject
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/bills/bills_api.md#reject

	Отклонить верификацию карты
*/
func (m *ModBills) RejectBillHandler(c *gin.Context) {
	req := &models.RejectBill{}
	if err := c.ShouldBindJSON(req); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if err := req.Validation(); err != nil {
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
		"created_at": time.Now().UTC().Format(core.DateStandart),
	})
	if err != nil {
		m.modlog(err)
	}

	if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
		m.modlog(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

/*
	@Method DELETE
	@Path admin/bill
	@Type PRIVATE
	@Documentation https://github.com/exchanger-bot/docs/blob/main/admin/bills/bills_api.md#create

	Создать запись в таблице `bills`

	# TESTED
*/
func (m *ModBills) CreateBillHandler(c *gin.Context) {
	r := &models.Bill{}
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, AppError.ErrInvalidBody)
		return
	}

	if obj := m.responser.RecordHandler(c, r, r.Validation()); obj != nil {
		if reflect.TypeOf(obj) != reflect.TypeOf(&models.Bill{}) {
			return
		}

		if err := m.responser.CreateRecordResponse(c, m.store.AdminPanel().Bills(), obj); err == nil {
			payload, err := json.Marshal(map[string]interface{}{
				"to": map[string]interface{}{
					"chat_id": r.ChatID,
				},
				"message": map[string]interface{}{
					"type": "confirmation_successful",
					"text": fmt.Sprintf("🟢 Успех! 🟢\n\nКарта `%s` успешно верифицированa.", r.Bill),
				},
				"created_at": time.Now().UTC().Format(core.DateStandart),
			})
			if err != nil {
				m.modlog(err)
				return
			}

			if err := m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
				m.modlog(err)
				return
			}
		}

		return
	}
}
