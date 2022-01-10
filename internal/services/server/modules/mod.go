package modules

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	AppType "github.com/gefion-tech/tg-exchanger-server/internal/core/types"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/guard"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/middleware"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/modules/bills"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/modules/exchanger"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/modules/logs"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/modules/message"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/modules/notification"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/server/modules/user"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type ServerModules struct {
	exMod     exchanger.ModExchangerI
	notifyMod notification.ModNotificationI
	userMod   user.ModUsersI
	msgMod    message.ModMessageI
	billsMod  bills.ModBillsI
	logsMod   logs.ModLogsI
}

type ServerModulesI interface {
	ModulesConfigure(router *gin.RouterGroup, g guard.GuardI, mdl middleware.MiddlewareI)
}

func InitServerModules(
	store db.SQLStoreI,
	redis *redisstore.AppRedisDictionaries,
	nsq nsqstore.NsqI,
	cnf *config.Config,
	logger utils.LoggerI,
	responser utils.ResponserI,
) ServerModulesI {
	return &ServerModules{
		exMod:     exchanger.InitModExchanger(store, redis, nsq, cnf, responser, logger),
		notifyMod: notification.InitModNotification(store, redis, nsq, cnf, logger, responser),
		msgMod:    message.InitModMessage(store, redis, nsq, cnf, responser, logger),
		userMod:   user.InitModUsers(store, redis, nsq, cnf, responser, logger),
		billsMod:  bills.InitModBills(store, redis, nsq, cnf, responser, logger),
		logsMod:   logs.InitModLogs(store.AdminPanel().Logs(), cnf, responser),
	}
}

func (m *ServerModules) ModulesConfigure(router *gin.RouterGroup, g guard.GuardI, mdl middleware.MiddlewareI) {
	// base
	{
		router.POST("/bot/user/registration",
			m.userMod.UserInBotRegistrationHandler,
		)
	}

	// bot bill
	{
		router.GET("/bot/user/bill/:id",
			m.billsMod.GetBillHandler,
		)
		router.GET("/bot/user/:chat_id/bills",
			m.billsMod.GetAllBillsHandler,
		)
		router.DELETE("/bot/user/:chat_id/bill/:id",
			m.billsMod.DeleteBillHandler,
		)
	}

	// bill
	{
		router.POST("/admin/bill",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.Logger(AppType.ResourceBill, AppType.ResourceCreate),
			m.billsMod.CreateBillHandler,
		)
		router.POST("/admin/bill/reject",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.Logger(AppType.ResourceBill, AppType.ResourceReject),
			m.billsMod.RejectBillHandler,
		)
	}

	// registration|auth
	{
		router.POST("/admin/registration/code",
			m.userMod.UserGenerateCodeHandler,
		)
		router.POST(
			"/admin/registration",
			m.userMod.UserInAdminRegistrationHandler,
		)
		router.POST("/admin/auth",
			m.userMod.UserInAdminAuthHandler,
		)
		router.POST("/admin/token/refresh",
			m.userMod.UserRefreshToken,
		)
		router.POST("/admin/logout",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.userMod.LogoutHandler,
		)
	}

	// message
	{
		router.POST("/admin/message",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.msgMod.CreateNewMessageHandler,
		)
		router.PUT("/admin/message/:id",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.Logger(AppType.ResourceMessage, AppType.ResourceUpdate),
			m.msgMod.UpdateBotMessageHandler,
		)
		router.DELETE("/admin/message/:id",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.Logger(AppType.ResourceMessage, AppType.ResourceDelete),
			m.msgMod.DeleteBotMessageHandler,
		)
		router.GET("/admin/message/:connector",
			m.msgMod.GetMessageHandler,
		)
		router.GET("/admin/messages",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.msgMod.GetMessagesSelectionHandler,
		)
	}

	// notification
	{
		router.POST("/admin/notification",
			m.notifyMod.CreateNotificationHandler,
		)
		router.PUT("/admin/notification/:id",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.notifyMod.UpdateNotificationStatusHandler,
		)
		router.DELETE("/admin/notification/:id",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.Logger(AppType.ResourceNotify, AppType.ResourceDelete),
			m.notifyMod.DeleteNotificationHandler,
		)
		router.GET("/admin/notifications",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.notifyMod.GetNotificationsSelectionHandler,
		)
		router.GET("/admin/notifications/check",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.notifyMod.NewNotificationsCheckHandler,
		)
	}

	// exchanger
	{
		router.POST("/admin/exchanger",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.exMod.CreateExchangerHandler,
		)
		router.PUT("/admin/exchanger/:id",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.Logger(AppType.ResourceExchange, AppType.ResourceUpdate),
			m.exMod.UpdateExchangerHandler,
		)
		router.DELETE("/admin/exchanger/:id",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.Logger(AppType.ResourceExchange, AppType.ResourceDelete),
			m.exMod.DeleteExchangerHandler,
		)
		router.GET("/admin/exchanger/:name",
			m.exMod.GetExchangerByNameHandler,
		)
		router.GET("/admin/exchanger/document",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.exMod.GetExchangerDocumentHandler,
		)
		router.GET("/admin/exchangers",
			g.AuthTokenValidation(),
			g.IsAuth(),
			m.exMod.GetExchangersSelectionHandler,
		)
	}

	// log
	{
		router.POST("/log",
			m.logsMod.CreateLogRecordHandler,
		)
		router.DELETE("/log/:id",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.IsAdmin(),
			m.logsMod.DeleteLogRecordHandler,
		)
		router.GET("/logs",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.IsAdmin(),
			m.logsMod.GetLogRecordsSelectionHandler,
		)
		router.DELETE("/logs",
			g.AuthTokenValidation(),
			g.IsAuth(),
			g.IsAdmin(),
			m.logsMod.DeleteLogRecordsSelectionHandler,
		)
	}
}
