package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

type ManagerRepository struct {
	store                 *sql.DB
	botMessagesRepository *BotMessagesRepository
}

/*
	==========================================================================================
	КОНСТРУКТОРЫ ВЛОЖЕННЫХ СТРУКТУР
	==========================================================================================
*/

func (r *ManagerRepository) BotMessages() db.BotMessagesRepository {
	if r.botMessagesRepository != nil {
		return r.botMessagesRepository
	}

	r.botMessagesRepository = &BotMessagesRepository{
		store: r.store,
	}

	return r.botMessagesRepository
}
