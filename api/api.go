package api

import (
	"bytes"
	"encoding/json"
	"github.com/KnightHacks/knighthacks_cli/config"
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
	query, err := BuildQuery("query Login($code: String!, $provider: Provider!) {login(code: $code, provider: $provider) {accountExists user{id} accessToken refreshToken encryptedOAuthAccessToken}}", map[string]any{"provider": provider, "code": code})
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

func (a *Api) Register(provider string, encryptedOAuthAccessToken string, user model.NewUser) (*model.RegistrationPayload, error) {
	query, err := BuildQuery(
		"mutation Register($provider: Provider!, $encryptedOAuthAccessToken: String!, $input: NewUser!) {register(encryptedOAuthAccessToken: $encryptedOAuthAccessToken, input: $input, provider: $provider) {user{id} accessToken refreshToken}}",
		map[string]any{
			"provider":                  provider,
			"encryptedOAuthAccessToken": encryptedOAuthAccessToken,
			"input":                     user,
		},
	)
	if err != nil {
		return nil, err
	}
	response, err := a.Client.Post(a.Endpoint, "application/json", bytes.NewReader(query))
	if err != nil {
		return nil, err
	}
	var parsedResponse struct {
		Data struct {
			RegistrationPayload model.RegistrationPayload `json:"register"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}
	err = ParseResponse(response.Body, &parsedResponse)
	if err != nil {
		return nil, err
	}
	HandleGraphQLErrors(parsedResponse.Errors)
	return &parsedResponse.Data.RegistrationPayload, nil
}

func (a *Api) Me(c *config.Config) (*model.User, error) {
	query, err := BuildQuery(
		"query Me {me { id firstName lastName email phoneNumber role age pronouns { subjective objective }}}",
		map[string]any{},
	)
	if err != nil {
		return nil, err
	}

	var parsedResponse struct {
		Data struct {
			User model.User `json:"me"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}
	err = MakeRequestWithHeaders(a, query, map[string]string{"Content-Type": "application/json", "authorization": c.Auth.Tokens.Access}, &parsedResponse)

	if err != nil {
		return nil, err
	}
	HandleGraphQLErrors(parsedResponse.Errors)
	return &parsedResponse.Data.User, nil
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

func MakeRequestWithHeaders[T interface{}](a *Api, body []byte, headers map[string]string, responseStruct *T) error {
	request, err := http.NewRequest("POST", a.Endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response, err := a.Client.Do(request)
	if err != nil {
		return err
	}
	return ParseResponse(response.Body, responseStruct)
}
