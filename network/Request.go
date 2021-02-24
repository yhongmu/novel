package network

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var Client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func GetSearchRequest(requestURL string, fun func(r io.Reader) (interface{}, error)) (interface{}, error) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 设定关闭响应体
	defer resp.Body.Close()
	return fun(resp.Body)
}

func GetRequest(requestURL string, urlValues url.Values) (string, error) {
	reqBody := urlValues.Encode()
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.URL.RawQuery = reqBody
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 设定关闭响应体
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func PostRequest(requestURL string, urlValues url.Values) (string, error) {
	reqBody := urlValues.Encode()
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.URL.RawQuery = reqBody
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 设定关闭响应体
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
