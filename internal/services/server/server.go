package server

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/app/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/private"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/public"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.SQLStoreI
	Router *gin.Engine
	config *config.Config

	// guard          guard.IGuard
	public_routes  public.PublicRoutesI
	private_routes private.PrivateRoutesI
}

type ServerI interface {
	Run() error
}

func Init(s db.SQLStoreI, r db.RedisStoreI, c *config.Config) ServerI {
	return root(s, r, c)
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
	s.private_routes.ConfigurePrivateRouter(v1)
}
func root(s db.SQLStoreI, r db.RedisStoreI, c *config.Config) *Server {
	// Инициализация роутера
	router := gin.New()

	// Инициализация модуля публичных маршрутов
	pub := public.Init(s, r, router, &c.Secrets)

	// Инициализация модуля приватных маршрутов
	prv := private.Init(s, r, router, &c.Secrets)

	server := &Server{
		store:  s,
		Router: router,
		config: c,
		// guard:          guard,
		public_routes:  pub,
		private_routes: prv,
	}

	gin.SetMode(gin.DebugMode)
	gin.ForceConsoleColor()

	server.configure()
	return server
}
