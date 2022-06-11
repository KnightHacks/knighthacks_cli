package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Api struct {
	Client   *http.Client
	Endpoint string
}

func (a Api) GetAuthRedirectLink(provider string) (string, error) {
	b, err := BuildQuery("query Query($provider: Provider!) {getAuthRedirectLink(provider: $provider)}", map[string]string{"provider": provider})
	if err != nil {
		return "", err
	}
	log.Printf("request=%s\n", string(b))
	response, err := a.Client.Post(a.Endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	log.Printf("response=%v\n", response)
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	log.Printf("repsonse=%s\n", all)

	var parsedResponse struct {
		Data struct {
			AuthRedirectLink string `json:"getAuthRedirectLink"`
		} `json:"data"`
	}
	err = json.Unmarshal(all, &parsedResponse)
	if err != nil {
		return "", err
	}

	return parsedResponse.Data.AuthRedirectLink, nil
}

func (a Api) Login(provider string, code string) {

}

func Query()

func BuildQuery(query string, variables map[string]string) ([]byte, error) {
	return json.Marshal(struct {
		Query     string            `json:"query"`
		Variables map[string]string `json:"variables"`
	}{query, variables})
}
