package exchanger

import (
	"database/sql"
	"encoding/xml"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/app/errors"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

/*
	@Method POST
	@Path admin/exchangers/:name
	@Type PUBLIC
	@Documentation

	Получение одной записи из таблицы `exchangers`

	# TESTED
*/
func (m *ModExchanger) GetExchangerByNameHandler(c *gin.Context) {
	// Операция получения записи из БД
	e, err := m.store.Manager().Exchanger().GetByName(&models.Exchanger{Name: c.Param("name")})
	switch err {
	case nil:
		c.JSON(http.StatusOK, e)
		return
	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound, errors.ErrRecordNotFound)
		return

	default:
		tools.ServErr(c, http.StatusUnprocessableEntity, err)
		return
	}
}

func (m *ModExchanger) GetExchangerDocumentHandler(c *gin.Context) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://1obmen.net/request-exportxml.xml")
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		tools.ServErr(c, http.StatusInternalServerError, err)
	}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	data := models.OneObmen{}
	if err := xml.Unmarshal(res.Body(), &data); err != nil {
		tools.ServErr(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"file": tools.OneObmenDocumentGenerate(&data, m.cnf.Server.Tmp),
	})
}
