package sqlstore

import (
	"database/sql"
)

type ExchangerRepository struct {
	store *sql.DB
}

/*
	==========================================================================================
	КОНЕЧНЫЕ МЕТОДЫ ТЕКУЩЕЙ СТРУКТУРЫ
	==========================================================================================
*/

// func (r *ExchangerRepository) Create() (*models.Exchanger, error) {

// }
