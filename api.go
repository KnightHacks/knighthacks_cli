package main

import (
	"bytes"
	"encoding/json"
	"github.com/KnightHacks/knighthacks_cli/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GraphQLError struct {
	Message string `json:"message"`
}

type Api struct {
	Client   *http.Client
	Endpoint string
}

func (a *Api) GetAuthRedirectLink(provider string) (string, error) {
	query, err := BuildQuery("query Login($provider: Provider!) {getAuthRedirectLink(provider: $provider)}", map[string]any{"provider": provider})
	if err != nil {
		return "", err
	}
	response, err := a.Client.Post(a.Endpoint, "application/json", bytes.NewReader(query))
	if err != nil {
		return "", err
	}

	var parsedResponse struct {
		Data struct {
			AuthRedirectLink string `json:"getAuthRedirectLink"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}

	err = ParseResponse(response.Body, &parsedResponse)
	if err != nil {
		return "", err
	}
	HandleGraphQLErrors(parsedResponse.Errors)
	return parsedResponse.Data.AuthRedirectLink, nil
}

func (a *Api) Login(provider string, code string) (*model.LoginPayload, error) {
	query, err := BuildQuery("query Login($code: String!, $provider: Provider!) {login(code: $code, provider: $provider) {accountExists user{id}}}", map[string]any{"provider": provider, "code": code})
	if err != nil {
		return nil, err
	}
	response, err := a.Client.Post(a.Endpoint, "application/json", bytes.NewReader(query))
	if err != nil {
		return nil, err
	}
	var parsedResponse struct {
		Data struct {
			Login model.LoginPayload `json:"login"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}
	err = ParseResponse(response.Body, &parsedResponse)
	if err != nil {
		return nil, err
	}
	HandleGraphQLErrors(parsedResponse.Errors)
	return &parsedResponse.Data.Login, nil
}

func (a *Api) Register(provider string, code string, user model.NewUser) (string, error) {
	query, err := BuildQuery(
		"mutation Register($provider: Provider!, $code: String!, $input: NewUser!) {register(code: $code, input: $input, provider: $provider) {id}}",
		map[string]any{
			"provider": provider,
			"code":     code,
			"input":    user,
		},
	)
	if err != nil {
		return "", err
	}
	response, err := a.Client.Post(a.Endpoint, "application/json", bytes.NewReader(query))
	if err != nil {
		return "", err
	}
	var parsedResponse struct {
		Data struct {
			Register struct {
				Id string `json:"id"`
			} `json:"register"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}
	err = ParseResponse(response.Body, &parsedResponse)
	if err != nil {
		return "", err
	}
	HandleGraphQLErrors(parsedResponse.Errors)
	return parsedResponse.Data.Register.Id, nil
}

func HandleGraphQLErrors(errs []GraphQLError) {
	if len(errs) > 0 {
		log.Println("The following errors occurred when attempting to register an account: ")
		for _, elem := range errs {
			log.Printf(elem.Message)
		}
		os.Exit(1)
	}
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
