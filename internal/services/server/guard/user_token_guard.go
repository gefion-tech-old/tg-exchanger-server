package guard

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenCtx int8

const CtxKeyToken TokenCtx = iota

/*
	Метод middleware для защиты приватных маршрутов
	В методе проверятеся, действителен ли токен или
	истек ли срок его действия.
*/
func (g *Guard) AuthTokenValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := g.validationToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
	Метод извлечения JWT токена из заголовка запроса
*/
func (g *Guard) extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

/*
	Метод верификации JWT токена
*/
func (g *Guard) verifyToken(r *http.Request) (*jwt.Token, error) {
	// Извлекаю токен, на выходе получаю просто строку
	tokenStr := g.extractToken(r)

	// Извлекаю токен в виде структуры
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Проверяю соответствие подписи токена с методом SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(g.secrets.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

/*
	Метод проверки токена на действительность,
	т.е. не истек ли его срок действия
*/
func (g *Guard) validationToken(r *http.Request) error {
	token, err := g.verifyToken(r)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

/*
	Метод извлечения метаданных токена
*/
func (g *Guard) extractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := g.verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		// Извлекаю chat_id из полезной нагрузки токена
		chatID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["chat_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &models.AccessDetails{
			AccessUuid: accessUuid,
			ChatID:     chatID,
			Username:   claims["username"].(string),
		}, nil
	}

	return nil, err
}
