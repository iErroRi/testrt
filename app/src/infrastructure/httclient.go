package infrastructure

import (
	"encoding/json"
	"net/http"
	"time"
)

func NewHttpClient() *HttpClient {
	httpClient := new(HttpClient)
	httpClient.instance = &http.Client{Timeout: 10 * time.Second}
	return httpClient
}

type HttpClient struct {
	instance *http.Client
}

func (httpClient *HttpClient) GetJson(url string, target interface{}) error {
	r, err := httpClient.instance.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
