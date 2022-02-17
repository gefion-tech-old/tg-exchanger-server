package directions

import (
	"github.com/gefion-tech/tg-exchanger-server/internal/config"
	"github.com/gefion-tech/tg-exchanger-server/internal/services/db"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type ModDirections struct {
	repository db.AdminPanelRepository
	cfg        *config.Config

	responser utils.ResponserI
}

type ModDirectionsI interface {
	CreateDirectionHandler(c *gin.Context)
}

func InitModDirections(
	r db.AdminPanelRepository,
	cfg *config.Config,
	responser utils.ResponserI,
) ModDirectionsI {
	return &ModDirections{
		repository: r,
		cfg:        cfg,
		responser:  responser,
	}
}
