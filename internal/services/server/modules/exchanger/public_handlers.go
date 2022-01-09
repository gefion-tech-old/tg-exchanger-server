package exchanger

import (
	"encoding/xml"
	"net/http"

	"github.com/gefion-tech/tg-exchanger-server/internal/models"
	"github.com/gefion-tech/tg-exchanger-server/internal/utils"
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
	r := &models.Exchanger{Name: c.Param("name")}
	m.responser.RecordResponse(c, r, m.store.AdminPanel().Exchanger().GetByName(r))
}

func (m *ModExchanger) GetExchangerDocumentHandler(c *gin.Context) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.SetRequestURI("https://1obmen.net/request-exportxml.xml")
	res := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, res); err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
	}

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	data := models.OneObmen{}
	if err := xml.Unmarshal(res.Body(), &data); err != nil {
		m.responser.Error(c, http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"file": utils.OneObmenDocumentGenerate(&data, m.cnf.Server.Tmp),
	})
}
