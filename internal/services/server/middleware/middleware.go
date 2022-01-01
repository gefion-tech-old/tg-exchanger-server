package middleware

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	logger utils.LoggerI
}

type MiddlewareI interface {
	CORSMiddleware() gin.HandlerFunc
}

func InitMiddleware(l utils.LoggerI) MiddlewareI {
	return &Middleware{
		logger: l,
	}
}
