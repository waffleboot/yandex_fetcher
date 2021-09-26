package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseYandexURL = "https://yandex.ru/search/touch/?service=www.yandex&ui=webmobileapp.yandex&numdoc=50&lr=213&p=0&text=%s"

func NewHttpClient() func(string) ([]byte, error) {
	var c http.Client
	return func(search string) ([]byte, error) {
		url := fmt.Sprintf(baseYandexURL, url.QueryEscape(search))
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}
