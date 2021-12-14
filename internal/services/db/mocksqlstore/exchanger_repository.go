package mocksqlstore

import "github.com/gefion-tech/tg-exchanger-server/internal/models"

type ExchangerRepository struct {
	exchangers map[uint]*models.Exchanger
}
