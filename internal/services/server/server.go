package server

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
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
	Run() error
}

func Init(s db.SQLStoreI, nsq nsqstore.NsqI, r *redisstore.AppRedisDictionaries, l utils.LoggerI, c *config.Config) ServerI {
	return root(s, nsq, r, l, c)
}

// Метод запуска сервера
func (s *Server) Run() error {
	return s.Router.Run(s.config.Server.Port)
}

// Метод общей конфигурации сервера
func (s *Server) configure() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")

	// Подключение всех модулей
	s.mods.ModulesConfigure(v1, s.guard, s.middleware)
}

func root(s db.SQLStoreI, nsq nsqstore.NsqI, r *redisstore.AppRedisDictionaries, l utils.LoggerI, c *config.Config) *Server {
	// Инициализация роутера
	router := gin.New()

	guard := guard.Init(r, &c.Secrets)
	m := middleware.InitMiddleware()

	router.Use(m.CORSMiddleware())

	server := &Server{
		store:      s,
		Router:     router,
		config:     c,
		guard:      guard,
		middleware: m,
		mods:       modules.InitServerModules(s, r, nsq, c, l, utils.InitResponser()),
	}

	gin.ForceConsoleColor()
	server.Router.Use(gin.Logger())

	server.configure()
	return server
}
