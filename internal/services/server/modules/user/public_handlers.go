package user

import (
	"database/sql"
	"encoding/json"
	_errors "errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/app/static"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
)

/*
	@Method POST
	@Path /bot/user/registration
	@Type PUBLIC
	@Documentation https://github.com/gefion-tech/tg-exchanger-server#registration-in-bot

	Регистрация человека как пользователя бота. При валидных данных создается
	запись в БД в таблице `users`.

	# TESTED
*/
func (m *ModUsers) UserInBotRegistrationHandler(c *gin.Context) {
	r := &models.User{}

	// Парсинг входящего тела запроса
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Выполнение операции с БД
	m.responser.NewRecordResponse(c, r, m.store.User().Create(r))
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
func (m *ModUsers) UserGenerateCodeHandler(c *gin.Context) {
	req := &models.UserFromAdminRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидирую полученные данные
	m.responser.Error(c, http.StatusUnprocessableEntity, req.UserFromAdminRequestValidation(m.cnf.Users))

	u, err := m.store.User().FindByUsername(req.Username)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			m.responser.Error(c, http.StatusNotFound, errors.ErrNotRegistered)
			return
		default:
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}
	}

	// Генерирую код подтверждения
	code := tools.VerificationCode(req.Testing)

	// Формирую сообщение и отправляю его в очередь
	msg := map[string]interface{}{
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

	payload, err := json.Marshal(msg)
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Хеширую пароль
	hash, err := tools.EncryptString(req.Password)
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Генерирую объект для записи в Redis
	b, err := json.Marshal(map[string]interface{}{
		"username": req.Username,
		"hash":     hash,
	})
	if err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Записываю в Redis
	m.responser.Error(c, http.StatusInternalServerError,
		m.redis.Registration.SaveVerificationCode(code, b),
	)

	// Отправляю сообщение в NSQ
	m.responser.Error(c, http.StatusInternalServerError,
		m.nsq.Publish(nsqstore.TOPIC__MESSAGES, payload),
	)

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
func (m *ModUsers) UserInAdminRegistrationHandler(c *gin.Context) {
	r := &models.UserCodeRequest{}

	// Парсинг входящего тела запроса
	if err := c.ShouldBindJSON(r); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Валидация
	m.responser.Error(c, http.StatusUnprocessableEntity, r.UserCodeRequestValidation())

	// Ищу данные по этому коду в Redis
	data, err := m.redis.Registration.FetchVerificationCode(int(r.Code))
	if err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity,
			_errors.New("activation period for this code has expired"),
		)
		return
	}

	u := &models.User{}
	if err := json.Unmarshal([]byte(data), u); err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

	// Назначение роли пользователю
	if r := tools.RoleDefine(u.Username, m.cnf.Users); r != static.S__ROLE__USER {
		u.Role = r
	} else {
		m.responser.Error(c, http.StatusForbidden, errors.ErrNotEnoughRights)
		return
	}

	m.responser.NewRecordResponse(c, u, m.store.User().RegisterInAdminPanel(u))
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
func (m *ModUsers) UserInAdminAuthHandler(c *gin.Context) {
	req := &models.UserFromAdminRequest{}

	if err := c.ShouldBindJSON(req); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	// Ищу пользователя в БД
	u, err := m.store.User().FindByUsername(req.Username)
	switch err {
	case nil:
		if u.Hash != nil && tools.ComparePassword(*u.Hash, req.Password) {
			// Генерирую сборку токенов и сопутствующих деталей
			td, err := m.createToken(u.ChatID, u.Username, u.Role)
			if err != nil {
				m.responser.Error(c, http.StatusInternalServerError, err)
				return
			}

			// Аутентифицирую пользователя
			if err := m.createAuth(u.ChatID, td); err != nil {
				m.responser.Error(c, http.StatusInternalServerError, err)
			}

			c.JSON(http.StatusOK, gin.H{
				"access_token":  td.AccessToken,
				"refresh_token": td.RefreshToken,
			})
			return
		}

		m.responser.Error(c, http.StatusNotFound, errors.ErrNotRegistered)
		return

	case sql.ErrNoRows:
		m.responser.Error(c, http.StatusNotFound, errors.ErrNotRegistered)
		return
	default:
		m.responser.Error(c, http.StatusInternalServerError, err)
		return
	}

}

/*
	@Method POST
	@Path /token/refresh
	@Type PUBLIC

	Метод обновления для токена доступа для доступа
	к приватным маршрутам.

*/
func (m *ModUsers) UserRefreshToken(c *gin.Context) {
	// Обрабатываю тело запроса пытаясь получить refresh токен
	mapToken := map[string]string{}

	if err := c.ShouldBindJSON(&mapToken); err != nil {
		m.responser.Error(c, http.StatusUnprocessableEntity, errors.ErrInvalidBody)
		return
	}

	refreshToken := mapToken["refresh_token"]

	// Верификация refresh токена
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Проверяю соответствие подписи токена с методом SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.cnf.Secrets.RefreshSecret), nil
	})

	// Если возникла ошибка, значит токен просрочен
	if err != nil {
		m.responser.Error(c, http.StatusUnauthorized, _errors.New("refresh token expired"))
		return
	}

	// Проверка валидности токена
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		m.responser.Error(c, http.StatusUnauthorized, err)
		return
	}

	// Если токен валиден получаю его uuid
	claims, ok := token.Claims.(jwt.MapClaims) // проверка на соответствие с MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) // конвертация интерфейса в строку
		if !ok {
			m.responser.Error(c, http.StatusUnprocessableEntity, err)
			return
		}

		// Извлекаю chat_id из полезной нагрузки токена
		chatID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["chat_id"]), 0, 64)
		if err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}

		// Удаляю предыдущий refresh токен
		deleted, err := m.redis.Auth.DeleteAuth(refreshUuid)
		if err != nil || deleted == 0 {
			m.responser.Error(c, http.StatusUnauthorized, _errors.New("unauthorized"))
			return
		}

		// Создание новой пары токенов
		ts, err := m.createToken(chatID, claims["username"].(string), int(claims["role"].(float64)))
		if err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}

		// Сохранение метаданных токенов в redis
		if err := m.createAuth(chatID, ts); err != nil {
			m.responser.Error(c, http.StatusInternalServerError, err)
			return
		}

		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		c.JSON(http.StatusCreated, tokens)
	} else {
		m.responser.Error(c, http.StatusUnauthorized, _errors.New("refresh expired"))
	}
}

/*
	==========================================================================================
	ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
	==========================================================================================
*/

/*
	Метод генерации пользовательского набора токенов
	на основе данных о пользователе
*/
func (m *ModUsers) createToken(chatID int64, username string, role int) (*models.TokenDetails, error) {
	// Набор информации о пользовательских токенах и иж сроки действия
	td := &models.TokenDetails{}
	var err error

	/* Определение времени жизни для токенов */

	// Определяю время жизни в 15 МИНУТ для токена ДОСТУПА
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	// Создаю идентификатор для токена доступа
	td.AccessUuid = uuid.NewV4().String()

	// Определяю время жизни в 7 ДНЕЙ для токена ОБНОВЛЕНИЯ
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	/* Генерация токена доступа */

	// Создаю полезную нагрузку токена
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["chat_id"] = chatID
	atClaims["username"] = username
	atClaims["role"] = role
	atClaims["exp"] = td.AtExpires

	// Кодирую полезную нагрузку создавая ТОКЕН ДОСТУПА
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(m.cnf.Secrets.AccessSecret))
	if err != nil {
		return nil, err
	}

	/* Генерация токена обновления */
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["chat_id"] = chatID
	rtClaims["username"] = username
	rtClaims["role"] = role
	rtClaims["exp"] = td.RtExpires

	// Кодирую полезную нагрузку создавая ТОКЕН ОБНОВЛЕНИЯ
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(m.cnf.Secrets.RefreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

/*
	Метод сохранения метаданных JWT
*/
func (m *ModUsers) createAuth(chatID int64, td *models.TokenDetails) error {
	// Конвертация access_token из Unix формата в UTC
	at := time.Unix(td.AtExpires, 0)
	// Конвертация refresh_token из Unix формата в UTC
	rt := time.Unix(td.RtExpires, 0)

	now := time.Now()

	// Сохранение access_tokenа
	if errAccess := m.redis.Auth.SaveAuth(td.AccessUuid, chatID, at.Sub(now)); errAccess != nil {
		return errAccess
	}

	// Сохранение refresh_tokenа
	if errRefresh := m.redis.Auth.SaveAuth(td.RefreshUuid, chatID, rt.Sub(now)); errRefresh != nil {
		return errRefresh
	}

	return nil
}
