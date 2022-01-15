package server

import (
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/guard"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/middleware"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/modules"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.SQLStoreI
	Router *gin.Engine
	config *config.Config

	guard      guard.GuardI
	middleware middleware.MiddlewareI
	mods       modules.ServerModulesI
}

type ServerI interface {
	Create() *http.Server
}

func Init(s db.SQLStoreI, nsq nsqstore.NsqI, r *redisstore.AppRedisDictionaries, l utils.LoggerI, c *config.Config) ServerI {
	return root(s, nsq, r, l, c)
}

func (s *Server) Create() *http.Server {
	return &http.Server{
		Addr:    s.config.Services.Server.Port,
		Handler: s.Router,
	}
}

func root(s db.SQLStoreI, nsq nsqstore.NsqI, r *redisstore.AppRedisDictionaries, l utils.LoggerI, c *config.Config) *Server {
	// Инициализация роутера
	router := gin.New()
	responser := utils.InitResponser(l)

	guard := guard.Init(r, &c.Secrets, responser, l)
	m := middleware.InitMiddleware(l)

	router.Use(m.CORSMiddleware())

	server := &Server{
		store:      s,
		Router:     router,
		config:     c,
		guard:      guard,
		middleware: m,
		mods:       modules.InitServerModules(s, r, nsq, c, l, responser),
	}

	gin.ForceConsoleColor()
	server.Router.Use(gin.Logger())

	server.configure()
	return server
}

// Метод общей конфигурации сервера
func (s *Server) configure() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")

	// Подключение всех модулей
	s.mods.ModulesConfigure(v1, s.guard, s.middleware)
}
