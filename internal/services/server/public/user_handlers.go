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
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
)

/*
	@Method POST
	@Path /bot/registration
	@Type PUBLIC
	@Documentation https://github.com/gefion-tech/tg-exchanger-server#registration-in-bot

	Регистрация человека как пользователя бота. При валидных данных создается
	запись в БД в таблице `users`.

	# TESTED
*/
func (pr *PublicRoutes) userInBotRegistrationHandler(c *gin.Context) {
	req := &models.User{}

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
		c.JSON(http.StatusCreated, u)
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
	@Documentation https://github.com/gefion-tech/tg-exchanger-server#registration-in-admin-panel

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

	// Валидирую полученные данные
	if err := req.UserFromAdminRequestValidation(pr.users.Managers, pr.users.Developers); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	u, err := pr.store.User().FindByUsername(req.Username)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": _errors.New("user with this username is not registered in bot").Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// Генерирую код подтверждения
	code := tools.VerificationCode(req.Testing)

	// Формирую сообщение и отправляю его в очередь
	m := map[string]interface{}{
		"to": map[string]interface{}{
			"chat_id":  u.ChatID,
			"username": u.Username,
		},
		"message": map[string]interface{}{
			"type": "verification_code",
			"text": fmt.Sprintf("%d", code),
		},
		"created_at": time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
	}

	payload, err := json.Marshal(m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Хеширую пароль
	hash, err := tools.EncryptString(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Записываю в Redis
	if err := pr.redis.Registration.SaveVerificationCode(code, b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println(m["to"].(map[string]interface{})["chat_id"])

	// Отправляю сообщение в NSQ
	if err := pr.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

/*
	@Method POST
	@Path /admin/registration
	@Type PUBLIC
	@Documentation https://github.com/gefion-tech/tg-exchanger-server#registration-in-admin-panel

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
	data, err := pr.redis.Registration.FetchVerificationCode(int(req.Code))
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
		c.JSON(http.StatusCreated, user)
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
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

/*
	@Method POST
	@Path /admin/auth
	@Type PUBLIC
	@Documentation https://github.com/gefion-tech/tg-exchanger-server#auth-in-admin-panel

	В методе проверяется, есть ли в бд пользователь с переданным username.
	Если пользователь найден, смотрим есть ли у него hash пароль (если нет, значит он не зареган как менеджер)
	Если хеш найден и совпадает с переданным паролем, создаю пользовательскую сессию.

	# TESTED
*/
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
		fmt.Println(tools.EncryptString(req.Password))
		// fmt.Println(*u.Hash)
		if u.Hash != nil && tools.ComparePassword(*u.Hash, req.Password) {
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
		}

		c.JSON(http.StatusNotFound, gin.H{
			"error": _errors.New("user with this username or password is not registered as manager").Error(),
		})
		return

	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{
			"error": _errors.New("user with this username or password is not registered as manager").Error(),
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

}
