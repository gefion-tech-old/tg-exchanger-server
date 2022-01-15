package whitebit_plugin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var _ apiHelperI = (*apiHelper)(nil)

type apiHelper struct {
	PublicKey string
	SecretKey string
	BaseURL   string
}

type apiHelperI interface {
	SendRequest(requestURL string, data map[string]string) ([]byte, error)
}

func (api *apiHelper) SendRequest(requestURL string, data map[string]string) ([]byte, error) {
	// Если одноразовый номер похож на номер предыдущего запроса или меньше
	// его, будет получено сообщение об ошибке «слишком много запросов»

	nonce := int(time.Now().Unix()) // nonce — это число, которое всегда больше, чем номер предыдущего запроса
	data["request"] = requestURL
	data["nonce"] = strconv.Itoa(nonce)

	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	completeURL := api.BaseURL + requestURL

	// Расчет полезной нагрузки
	payload := base64.StdEncoding.EncodeToString(requestBody)

	// Вычисление подписи с использованием sha512
	h := hmac.New(sha512.New, []byte(api.SecretKey))
	h.Write([]byte(payload))
	signature := fmt.Sprintf("%x", h.Sum(nil))

	client := http.Client{}

	request, err := http.NewRequest("POST", completeURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-type", "application/json")
	request.Header.Set("X-TXC-APIKEY", api.PublicKey)
	request.Header.Set("X-TXC-PAYLOAD", payload)
	request.Header.Set("X-TXC-SIGNATURE", signature)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
