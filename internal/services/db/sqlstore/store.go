package sqlstore

import (
	"database/sql"

	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

// SQL хранилище
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func Init(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() db.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s.db,
	}

	return s.userRepository
}
