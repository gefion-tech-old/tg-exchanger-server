package mocksqlstore

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
)

// Mock SQL хранилище
type Store struct {
	userRepository *UserRepository
}

func Init() db.SQLStoreI {
	return &Store{}
}

func (s *Store) User() db.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		users: make(map[int64]*models.User),
	}

	return s.userRepository
}
