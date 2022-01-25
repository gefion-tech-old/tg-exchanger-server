package ma

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/plugins"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type ModMerchantAutoPayout struct {
	repository db.AdminPanelRepository
	redis      *redisstore.AppRedisDictionaries
	nsq        nsqstore.NsqI
	cfg        *config.Config
	pl         *plugins.AppPlugins

	responser utils.ResponserI
	logger    utils.LoggerI
}

type ModMerchantAutoPayoutI interface {
	CreateMerchantAutopayoutHandler(c *gin.Context)
	UpdateMerchantAutopayoutHandler(c *gin.Context)
	DeleteMerchantAutopayoutHandler(c *gin.Context)
	GetMerchantAutopayoutHandler(c *gin.Context)
	GetMerchantAutopayoutSelectionHandler(c *gin.Context)

	CreateNewAdressHandler(c *gin.Context)
}

func InitModMerchantAutoPayout(
	rep db.AdminPanelRepository,
	redis *redisstore.AppRedisDictionaries,
	nsq nsqstore.NsqI,
	cfg *config.Config,
	pl *plugins.AppPlugins,
	responser utils.ResponserI,
	l utils.LoggerI,
) ModMerchantAutoPayoutI {
	return &ModMerchantAutoPayout{
		repository: rep,
		redis:      redis,
		nsq:        nsq,
		cfg:        cfg,

		pl: pl,

		responser: responser,
		logger:    l,
	}
}
