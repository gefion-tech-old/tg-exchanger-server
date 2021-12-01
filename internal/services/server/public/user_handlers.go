package public

import (
	"database/sql"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path /registration
	@Type PUBLIC
	@Documentation

	Регистрация пользователя. При валидных данных создается
	запись в БД в таблице `users`.
*/
func (pr *PublicRoutes) userRegistrationHandler(c *gin.Context) {
	req := &models.UserRequest{}

	// Парсинг входящего тела запроса
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	// Выполнение операции с БД
	u, err := pr.store.User().Create(req)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, gin.H{
			"chat_id":    u.ChatID,
			"username":   u.Username,
			"hash":       u.Hash,
			"created_at": u.CreatedAt,
			"updated_at": u.UpdatedAt,
		})
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrAlreadyRegistered.Error(),
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}
