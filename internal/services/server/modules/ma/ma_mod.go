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
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cfg   *config.Config
	pl    *plugins.AppPlugins

	responser utils.ResponserI
	logger    utils.LoggerI
}

type ModMerchantAutoPayoutI interface {
	CreateNewAdressHandler(c *gin.Context)
}

func InitModMerchantAutoPayout(
	store db.SQLStoreI,
	redis *redisstore.AppRedisDictionaries,
	nsq nsqstore.NsqI,
	cfg *config.Config,
	pl *plugins.AppPlugins,
	responser utils.ResponserI,
	l utils.LoggerI,
) ModMerchantAutoPayoutI {
	return &ModMerchantAutoPayout{
		store: store,
		redis: redis,
		nsq:   nsq,
		cfg:   cfg,

		pl: pl,

		responser: responser,
		logger:    l,
	}
}
