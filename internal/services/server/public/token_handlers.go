package public

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

/*
	@Method POST
	@Path /token/refresh
	@Type PUBLIC

	Метод обновления для токена доступа для доступа
	к приватным маршрутам.

*/
func (pr *PublicRoutes) refreshToken(c *gin.Context) {
	// Обрабатываю тело запроса пытаясь получить refresh токен
	mapToken := map[string]string{}

	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	refreshToken := mapToken["refresh_token"]

	// Верификация refresh токена
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Проверяю соответствие подписи токена с методом SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(pr.secrets.RefreshSecret), nil
	})

	// Если возникла ошибка, значит токен просрочен
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token expired",
		})
		return
	}

	// Проверка валидности токена
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	// Если токен валиден получаю его uuid
	claims, ok := token.Claims.(jwt.MapClaims) // проверка на соответствие с MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) // конвертация интерфейса в строку
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err,
			})
			return
		}

		// Извлекаю chat_id из полезной нагрузки токена
		chatID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["chat_id"]), 0, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "error occurred",
			})
			return
		}

		// Удаляю предыдущий refresh токен
		// deleted, err := pr.redis.Auth.Del(refreshUuid).Result()
		// if err != nil || deleted == 0 {
		// 	c.JSON(http.StatusUnauthorized, gin.H{
		// 		"error": "unauthorized",
		// 	})
		// 	return
		// }

		deleted, err := pr.redis.Auth.DeleteAuth(refreshUuid)
		if err != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		// Создание новой пары токенов
		ts, err := pr.createToken(chatID, claims["username"].(string))
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "error occurred",
			})
			return
		}

		// Сохранение метаданных токенов в redis
		if err := pr.createAuth(chatID, ts); err != nil {
			c.JSON(http.StatusForbidden, err.Error())
			return
		}

		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
