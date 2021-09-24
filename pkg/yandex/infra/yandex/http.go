package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const baseYandexURL = "https://yandex.ru/search/touch/?service=www.yandex&ui=webmobileapp.yandex&numdoc=50&lr=213&p=0&text=%s"

type HTTPClient struct {
}

func NewHttpClient() HTTPClient {
	return HTTPClient{}
}

func (h HTTPClient) Supply(search string) ([]byte, error) {
	url := fmt.Sprintf(baseYandexURL, url.QueryEscape(search))
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
