package public

import (
	"strconv"
	"time"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/twinj/uuid"
)

type PublicRoutes struct {
	store   db.SQLStoreI
	redis   *redisstore.AppRedisDictionaries
	router  *gin.Engine
	secrets *config.SecretsConfig
	users   *config.UsersConfig
}

type PublicRoutesI interface {
	ConfigurePublicRouter(router *gin.RouterGroup)
}

// Конструктор модуля публичных маршрутов
func Init(store db.SQLStoreI, redis *redisstore.AppRedisDictionaries, router *gin.Engine, secrets *config.SecretsConfig, users *config.UsersConfig) PublicRoutesI {
	return &PublicRoutes{
		store:   store,
		redis:   redis,
		router:  router,
		secrets: secrets,
		users:   users,
	}
}

// Метод конфигуратор всех публичных маршрутов
func (pr *PublicRoutes) ConfigurePublicRouter(router *gin.RouterGroup) {
	router.POST("/token/refresh", pr.refreshToken)

	bot := router.Group("/bot")
	bot.POST("/registration", pr.userInBotRegistrationHandler)

	admin := router.Group("/admin")
	admin.POST("/registration/code", pr.userGenerateCodeHandler)
	admin.POST("/registration", pr.userInAdminRegistrationHandler)
	admin.POST("/auth", pr.userInAdminAuthHandler)
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
func (pr *PublicRoutes) createToken(chatID int64, username string) (*models.TokenDetails, error) {
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
	atClaims["exp"] = td.AtExpires

	// Кодирую полезную нагрузку создавая ТОКЕН ДОСТУПА
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(pr.secrets.AccessSecret))
	if err != nil {
		return nil, err
	}

	/* Генерация токена обновления */
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["chat_id"] = chatID
	rtClaims["username"] = username
	rtClaims["exp"] = td.RtExpires

	// Кодирую полезную нагрузку создавая ТОКЕН ОБНОВЛЕНИЯ
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(pr.secrets.RefreshSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

/*
	Метод сохранения метаданных JWT
*/
func (pr *PublicRoutes) createAuth(chatID int64, td *models.TokenDetails) error {
	// Конвертация access_token из Unix формата в UTC
	at := time.Unix(td.AtExpires, 0)
	// Конвертация refresh_token из Unix формата в UTC
	rt := time.Unix(td.RtExpires, 0)

	now := time.Now()

	// Сохранение access_tokenа
	if errAccess := pr.redis.Auth.Set(
		td.AccessUuid,
		strconv.Itoa(int(chatID)),
		at.Sub(now),
	).Err(); errAccess != nil {
		return errAccess
	}

	// Сохранение refresh_tokenа
	if errRefresh := pr.redis.Auth.Set(
		td.RefreshUuid,
		strconv.Itoa(int(chatID)),
		rt.Sub(now),
	).Err(); errRefresh != nil {
		return errRefresh
	}

	return nil
}
