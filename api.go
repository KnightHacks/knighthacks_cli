package main

import (
	"bytes"
	"encoding/json"
	"github.com/KnightHacks/knighthacks_cli/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Api struct {
	Client   *http.Client
	Endpoint string
}

func (a *Api) GetAuthRedirectLink(provider string) (string, error) {
	b, err := BuildQuery("query Login($provider: Provider!) {getAuthRedirectLink(provider: $provider)}", map[string]any{"provider": provider})
	if err != nil {
		return "", err
	}
	log.Printf("request=%s\n", string(b))
	response, err := a.Client.Post(a.Endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	var parsedResponse struct {
		Data struct {
			AuthRedirectLink string `json:"getAuthRedirectLink"`
		} `json:"data"`
	}

	err = ParseResponse(response.Body, &parsedResponse)
	return parsedResponse.Data.AuthRedirectLink, nil
}

func (a *Api) Login(provider string, code string) (*model.LoginPayload, error) {
	b, err := BuildQuery("query Login($code: String!, $provider: Provider!) {login(code: $code, provider: $provider) {accountExists user{id}}}", map[string]any{"provider": provider, "code": code})
	if err != nil {
		return nil, err
	}
	log.Printf("request=%s\n", string(b))
	response, err := a.Client.Post(a.Endpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var parsedResponse struct {
		Data struct {
			Login model.LoginPayload `json:"login"`
		} `json:"data"`
	}
	err = ParseResponse(response.Body, &parsedResponse)

	return &parsedResponse.Data.Login, nil
}

//ParseResponse
//	var response struct {
//		Data struct{} `json:"data,omitempty"`
//	}
// the contents of response should be what was previously set with the actual data inside the data struct
func ParseResponse[T interface{}](body io.ReadCloser, response *T) error {
	all, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	log.Printf("repsonse=%s\n", all)

	err = json.Unmarshal(all, response)
	return err

}

func BuildQuery(query string, variables map[string]any) ([]byte, error) {
	return json.Marshal(struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}{query, variables})
}
