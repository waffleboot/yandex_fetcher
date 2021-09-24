package application

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func Query(search string) error {
	url := fmt.Sprintf(baseYandexURL, url.QueryEscape(search))
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result := parseYandexResponse(data)
	if result.Error != nil {
		return result.Error
	}
	for _, v := range result.Items {
		log.Println(v.Host, v.Url)
	}
	return nil
}
