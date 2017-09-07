package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type WebAccess interface {
	Post(href_slug string, json_bytes []byte) ([]byte, error)
	Get(href_slug string) ([]byte, error)
}

type ConnectionParameters_t struct {
	BaseUrl    string
	Username   string
	Password   string
	MIQToken   string
	verify_ssl bool
	Group      string
}

func (params ConnectionParameters_t) Post(href_slug string, json_bytes []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", params.api_path(href_slug), bytes.NewBuffer(json_bytes))
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	return params.response(req)
}

func (params ConnectionParameters_t) Get(href_slug string) ([]byte, error) {
	fmt.Println("URL is ", params.api_path(href_slug))
	req, err := http.NewRequest("GET", params.api_path(href_slug), nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	return params.response(req)
}

func (params ConnectionParameters_t) response(req *http.Request) ([]byte, error) {
	if params.MIQToken != "" {
		req.Header.Set("X-Auth-Token", params.MIQToken)
	} else if params.Username != "" && params.Password != "" {
		req.SetBasicAuth(params.Username, params.Password)
	}

	if params.Group != "" {
		req.Header.Set("X-MIQ-Group", params.Group)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("client.Do: ", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body, nil
}

func (params ConnectionParameters_t) api_path(href_slug string) string {
	if strings.HasPrefix(href_slug, "http://") || strings.HasPrefix(href_slug, "https://") {
		return href_slug
	}
	return params.BaseUrl + href_slug
}
