package server

import (
	"os"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/guard"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/private"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/public"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	store  db.SQLStoreI
	Router *gin.Engine
	config *config.Config
	logger *logrus.Logger

	guard          guard.GuardI
	public_routes  public.PublicRoutesI
	private_routes private.PrivateRoutesI
}

type ServerI interface {
	Run() error
}

func Init(s db.SQLStoreI, nsq nsqstore.NsqI, r *redisstore.AppRedisDictionaries, c *config.Config) ServerI {
	return root(s, nsq, r, c)
}

// Метод запуска сервера
func (s *Server) Run() error {
	return s.Router.Run(s.config.Server.Port)
}

// Метод общей конфигурации сервера
func (s *Server) configure() {
	api := s.Router.Group("/api")
	v1 := api.Group("/v1")

	// Подключение публичных путей
	s.public_routes.ConfigurePublicRouter(v1)

	// Подключение приватных полей
	s.private_routes.ConfigurePrivateRouter(v1, s.guard)
}

func root(s db.SQLStoreI, nsq nsqstore.NsqI, r *redisstore.AppRedisDictionaries, c *config.Config) *Server {
	// Настрока логгера
	log := logrus.New()
	f, err := os.OpenFile("logs/test.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

	}

	log.SetOutput(f)
	log.SetLevel(logrus.ErrorLevel)

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Error("A group of walrus emerges from the ocean")

	// Инициализация роутера
	router := gin.New()

	// Инициализация охранников маршрутов
	guard := guard.Init(r, &c.Secrets)

	// Инициализация модуля публичных маршрутов
	pub := public.Init(s, nsq, r, router, &c.Secrets, &c.Users)

	// Инициализация модуля приватных маршрутов
	prv := private.Init(s, r, nsq, router, &c.Secrets, &c.Server, &c.Users)

	server := &Server{
		store:          s,
		Router:         router,
		config:         c,
		logger:         log,
		guard:          guard,
		public_routes:  pub,
		private_routes: prv,
	}

	gin.ForceConsoleColor()
	server.Router.Use(gin.Logger())

	server.configure()
	return server
}
