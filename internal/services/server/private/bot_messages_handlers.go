package private

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

/*
	@Method DELETE
	@Path admin/message/:id
	@Type PRIVATE
	@Documentation

	Удалить запись в таблице `bot_messages`

	# TESTED
*/
func (pr *PrivateRoutes) deleteBotMessageHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	msg, err := pr.store.Manager().BotMessages().Delete(&models.BotMessage{ID: uint(id)})
	switch err {
	case nil:
		c.JSON(http.StatusOK, msg)
		return

	case sql.ErrNoRows:
		tools.ServErr(c, http.StatusNotFound, errors.ErrRecordNotFound)
		return

	default:
		tools.ServErr(c, http.StatusInternalServerError, err)
		return
	}
}

/*
	@Method PUT
	@Path admin/message/:id
	@Type PRIVATE
	@Documentation

	Обновить запись в таблице `bot_messages`

	# TESTED
*/
func (pr *PrivateRoutes) updateBotMessageHandler(c *gin.Context) {
	req := &models.BotMessage{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	req.ID = uint(id)

	if err := req.UpdateBotMessageValidation(pr.users.Managers, pr.users.Developers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg, err := pr.store.Manager().BotMessages().Update(req)
	switch err {
	case nil:
		c.JSON(http.StatusOK, msg)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": errors.ErrRecordNotFound,
		})
		return

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

}

/*
	@Method GET
	@Path admin/messages
	@Type PRIVATE
	@Documentation

	Получение лимитированного объема записей из таблицы `bot_messages`

	# TESTED
*/
func (pr *PrivateRoutes) getAllBotMessageHandler(c *gin.Context) {
	errs, _ := errgroup.WithContext(c)

	cArrM := make(chan []*models.BotMessage)
	cCount := make(chan *int)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if err != nil {
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}

	// Подсчет кол-ва сообщений в таблице
	errs.Go(func() error {
		defer close(cCount)
		c, err := pr.store.Manager().BotMessages().Count()
		if err != nil {
			return err
		}

		cCount <- &c
		return nil
	})

	// Достаю из БД запрашиваемые записи
	errs.Go(func() error {
		defer close(cArrM)
		arrM, err := pr.store.Manager().BotMessages().GetSlice(page * limit)
		if err != nil {
			return err
		}

		cArrM <- arrM
		return nil
	})

	arrM := <-cArrM
	count := <-cCount

	if arrM == nil || count == nil {
		tools.ServErr(c, http.StatusInternalServerError, errs.Wait())
		return
	}

	fmt.Println(*count)

	c.JSON(http.StatusOK, gin.H{
		"limit":        limit,
		"current_page": page,
		"last_page":    math.Ceil(float64(*count) / float64(limit)),
		"data":         arrM[(tools.LowerThreshold(page, limit, *count)-1)*limit : tools.UpperThreshold(page, limit, *count)],
	})
}

/*
	@Method GET
	@Path admin/message/:connector
	@Type PRIVATE
	@Documentation

	Получить запись из таблицы `bot_messages`

	# TESTED
*/
func (pr *PrivateRoutes) getBotMessageHandler(c *gin.Context) {
	msg, err := pr.store.Manager().BotMessages().Get(&models.BotMessage{Connector: c.Param("connector")})
	switch err {
	case nil:
		c.JSON(http.StatusOK, msg)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": "message with current connector is not found",
		})
		return

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

}

/*
	@Method POST
	@Path admin/message
	@Type PRIVATE
	@Documentation

	Создать запись в таблице `bot_messages`

	# TESTED
*/
func (pr *PrivateRoutes) createNewBotMessageHandler(c *gin.Context) {
	req := &models.BotMessage{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	if err := req.CreateBotMessageValidation(pr.users.Managers, pr.users.Developers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg, err := pr.store.Manager().BotMessages().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, msg)
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusConflict, gin.H{
			"error": "message with current connector already created",
		})
		return

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}
