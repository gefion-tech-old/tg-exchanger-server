package middleware

import "github.com/gin-gonic/gin"

type Middleware struct{}

type MiddlewareI interface {
	CORSMiddleware() gin.HandlerFunc
	SetCors() gin.HandlerFunc
}

func InitMiddleware() MiddlewareI {
	return &Middleware{}
}
