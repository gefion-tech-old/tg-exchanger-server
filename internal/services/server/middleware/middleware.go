package middleware

import "github.com/gin-gonic/gin"

type Middleware struct{}

type MiddlewareI interface {
	SetCors() gin.HandlerFunc
}

func InitMiddleware() MiddlewareI {
	return &Middleware{}
}
