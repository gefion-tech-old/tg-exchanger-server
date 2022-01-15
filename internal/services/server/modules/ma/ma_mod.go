package ma

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/nsqstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db/redisstore"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
)

type ModMerchantAutoPayout struct {
	store db.SQLStoreI
	redis *redisstore.AppRedisDictionaries
	nsq   nsqstore.NsqI
	cnf   *config.Config

	responser utils.ResponserI
	logger    utils.LoggerI
}

type ModMerchantAutoPayoutI interface{}

func InitModMerchantAutoPayout(
	store db.SQLStoreI,
	redis *redisstore.AppRedisDictionaries,
	nsq nsqstore.NsqI,
	cfg *config.Config,
	responser utils.ResponserI,
	l utils.LoggerI,
) ModMerchantAutoPayoutI {
	return &ModMerchantAutoPayout{
		store: store,
		redis: redis,
		nsq:   nsq,
		cnf:   cfg,

		responser: responser,
		logger:    l,
	}
}
