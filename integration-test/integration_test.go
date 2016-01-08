package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var token string

func TestGetTopics(t *testing.T) {

	resp, err := getUnauthorized("/lecture-service/topics")
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp, err = getAuthorized("/lecture-service/topics")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	topics := make([]map[string]interface{}, 0)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&topics)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 2, len(topics))
	assert.Equal(t, topics[0]["name"].(string), "Grundlagen der Programmierung mit Java")
}

func getUnauthorized(path string) (*http.Response, error) {
	host := getHost()
	if token == "" {
		token = getToken()
	}
	req, err := http.NewRequest("GET", "http://"+host+path, nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{}
	return client.Do(req)
}

func getAuthorized(path string) (*http.Response, error) {
	host := getHost()
	if token == "" {
		token = getToken()
	}
	req, err := http.NewRequest("GET", "http://"+host+path, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := http.Client{}
	return client.Do(req)
}

func getToken() string {
	host := getHost()
	reader := strings.NewReader("client_id=user-web-client&client_secret=user-web-client-secret&grant_type=password&username=admin&password=admin")
	result, err := http.Post("http://"+host+"/authentication-service/oauth/token", "application/x-www-form-urlencoded", reader)
	if err != nil {
		panic(err)
	}
	defer result.Body.Close()
	if result.StatusCode != 200 {
		panic(Errors.New("expected statuscode 200 from authentication-service, but got: ", result.StatusCode, result.Status))
	}
	token := make(map[string]interface{})
	err = json.NewDecoder(result.Body).Decode(&token)
	if err != nil {
		panic(err)
	}
	return token["access_token"].(string)
}

func getHost() string {
	return os.Getenv("DH")
}
