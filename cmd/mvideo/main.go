package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if err := run("https://www.mvideo.ru/playstation-4327"); err != nil {
		log.Println(err)
	}
}

type MyJar struct {
	cookies []*http.Cookie
}

func (j *MyJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	// log.Printf("set cookies: %v %v\n", u, cookies)
	j.cookies = cookies
}

func (j *MyJar) Cookies(u *url.URL) []*http.Cookie {
	return j.cookies
}

func run(url string) error {
	out, err := os.Create("mvideo.html")
	if err != nil {
		return err
	}
	defer out.Close()
	client := http.Client{
		Jar: &MyJar{},
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("DNT", "1")
	// req.Header.Set("Host", "www.mvideo.ru")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "none")
	// req.Header.Set("Sec-Fetch-User", "?1")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:92.0) Gecko/20100101 Firefox/92.0")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		fmt.Printf("header: %v %v\n", k, v)
	}
	fmt.Printf("status code: %d\n", resp.StatusCode)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
