package public

import (
	"database/sql"
	"encoding/json"
	_errors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path /bot/registration
	@Type PUBLIC
	@Documentation

	Регистрация человека как пользователя бота. При валидных данных создается
	запись в БД в таблице `users`.

	# TESTED
*/
func (pr *PublicRoutes) userInBotRegistrationHandler(c *gin.Context) {
	req := &models.UserFromBotRequest{}

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

/*
	@Method POST
	@Path /admin/registration/code
	@Type PUBLIC
	@Documentation

	Пользователю с переданным username в ЛС будет отправлен код подтверждения
	который он должен будет ввести в окне на фронтенде. В ответ на запрос отдается
	код который был отправлен человеку в ЛС.

	# TESTED
*/
func (pr *PublicRoutes) userGenerateCodeHandler(c *gin.Context) {
	req := &models.UserFromAdminRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	// Генерирую код подтверждения
	code := tools.RandInt(100000, 999999)

	// Валидирую данные
	if err := req.UserFromAdminRequestValidation(pr.users.Managers, pr.users.Developers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Хеширую пароль
	hash, err := tools.EncryptString(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Генерирую объект для записи в Redis
	b, err := json.Marshal(map[string]interface{}{
		"username": req.Username,
		"hash":     hash,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Записываю в Redis
	if err := pr.redis.Registration.Set(fmt.Sprintf("%d", code), b, 30*time.Minute).Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code, // УБРАТЬ ЭТУ СТРОКУ В ПРОДЕ
	})
}

/*
	@Method POST
	@Path /admin/registration
	@Type PUBLIC
	@Documentation

	Регистрация человека как пользователя фвьин-панели. При валидных данных
	обновляется поле has в БД в таблице `users`.

	# TESTED
*/
func (pr *PublicRoutes) userInAdminRegistrationHandler(c *gin.Context) {
	req := &models.UserCodeRequest{}

	// Парсинг входящего тела запроса
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	// Валидация
	if err := req.UserCodeRequestValidation(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Ищу данные по этому коду в Redis
	data, err := pr.redis.Registration.Get(fmt.Sprintf("%d", req.Code)).Result()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": _errors.New("activation period for this code has expired").Error(),
		})
		return
	}

	u := models.User{}
	if err := json.Unmarshal([]byte(data), &u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := pr.store.User().RegisterAsManager(&u)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, gin.H{
			"chat_id":    user.ChatID,
			"username":   user.Username,
			"hash":       user.Hash,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrNotRegistered.Error(),
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (pr *PublicRoutes) userInAdminAuthHandler(c *gin.Context) {
	req := &models.UserFromAdminRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrInvalidBody.Error(),
		})
		return
	}

	// Ищу пользователя в БД
	u, err := pr.store.User().FindByUsername(req.Username)
	switch err {
	case nil:
		// Генерирую сборку токенов и сопутствующих деталей
		td, err := pr.createToken(u.ChatID, u.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Аутентифицирую пользователя
		if err := pr.createAuth(u.ChatID, td); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  td.AccessToken,
			"refresh_token": td.RefreshToken,
		})
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": errors.ErrNotRegistered.Error(),
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

}
