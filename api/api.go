package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/KnightHacks/knighthacks_cli/config"
	"github.com/KnightHacks/knighthacks_cli/model"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

type GraphQLError struct {
	Message string `json:"message"`
}

type Api struct {
	Client    *http.Client
	Endpoint  string
	DebugMode bool
}

func NewApi() (*Api, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Api{Client: &http.Client{Timeout: time.Second * 10, Jar: jar}}, nil
}

func (a *Api) GetAuthRedirectLink(provider string) (string, string, error) {
	query, err := BuildQuery("query Login($provider: Provider!) {getAuthRedirectLink(provider: $provider)}", map[string]any{"provider": provider})
	if err != nil {
		return "", "", err
	}

	var parsedResponse struct {
		Data struct {
			AuthRedirectLink string `json:"getAuthRedirectLink"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}

	response, err := MakeRequestWithHeaders(a, query, map[string]string{}, &parsedResponse, nil)
	if err != nil {
		return "", "", err
	}

	HandleGraphQLErrors(parsedResponse.Errors)
	cookies := response.Cookies()
	if a.DebugMode {
		fmt.Printf("cookies=%s\n", cookies)
	}
	var oauthState string
	for _, elem := range cookies {
		if elem.Name == "oauthstate" {
			oauthState = elem.Value
			break
		}
	}
	if len(oauthState) == 0 {
		return "", "", fmt.Errorf("unable to find oauthstate cookie")
	}
	oauthState, err = url.QueryUnescape(oauthState)
	if err != nil {
		return "", "", err
	}
	return parsedResponse.Data.AuthRedirectLink, oauthState, nil
}

func (a *Api) Login(provider string, code string, state string) (*model.LoginPayload, error) {
	query, err := BuildQuery(
		"query Login($code: String!, $state: String!, $provider: Provider!) {login(code: $code, provider: $provider, state: $state) {accountExists user{id} accessToken refreshToken encryptedOAuthAccessToken}}",
		map[string]any{"provider": provider, "code": code, "state": state},
	)
	if err != nil {
		return nil, err
	}
	var parsedResponse struct {
		Data struct {
			Login model.LoginPayload `json:"login"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}

	_, err = MakeRequestWithHeaders(a, query, map[string]string{}, &parsedResponse, map[string]string{"oauthstate": state})
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

	var parsedResponse struct {
		Data struct {
			RegistrationPayload model.RegistrationPayload `json:"register"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}

	_, err = MakeRequestWithHeaders(a, query, map[string]string{}, &parsedResponse, nil)
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
	_, err = MakeRequestWithHeaders(a, query, map[string]string{"authorization": c.Auth.Tokens.Access}, &parsedResponse, nil)

	if err != nil {
		return nil, err
	}
	HandleGraphQLErrors(parsedResponse.Errors)
	return &parsedResponse.Data.User, nil
}

func (a *Api) Delete(c *config.Config, id string) (bool, error) {
	query, err := BuildQuery(
		"mutation DeleteUser($userId: ID!) {deleteUser(id: $userId)}",
		map[string]any{"userId": id},
	)
	if err != nil {
		return false, err
	}

	var parsedResponse struct {
		Data struct {
			Deleted bool `json:"deleteUser"`
		} `json:"data"`
		Errors []GraphQLError `json:"errors"`
	}
	_, err = MakeRequestWithHeaders(a, query, map[string]string{"authorization": c.Auth.Tokens.Access}, &parsedResponse, nil)

	if err != nil {
		return false, err
	}
	HandleGraphQLErrors(parsedResponse.Errors)
	return parsedResponse.Data.Deleted, nil
}

func HandleGraphQLErrors(errs []GraphQLError) {
	if len(errs) > 0 {
		log.Println("The following errors occurred when attempting to handle your graphql query: ")
		for _, elem := range errs {
			log.Printf(elem.Message)
		}
		os.Exit(1)
	}
}

// ParseResponse
//
//	var response struct {
//		Data struct{} `json:"data,omitempty"`
//	}
//
// the contents of response should be what was previously set with the actual data inside the data struct
func ParseResponse[T interface{}](api *Api, body io.ReadCloser, response *T) error {
	all, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if api.DebugMode {
		log.Printf("repsonse=%s\n", all)
	}
	err = json.Unmarshal(all, response)
	return err
}

func BuildQuery(query string, variables map[string]any) ([]byte, error) {
	return json.Marshal(struct {
		Query     string         `json:"query"`
		Variables map[string]any `json:"variables"`
	}{query, variables})
}

func MakeRequestWithHeaders[T interface{}](a *Api, body []byte, headers map[string]string, responseStruct *T, cookies map[string]string) (*http.Response, error) {
	request, err := http.NewRequest("POST", a.Endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	if cookies != nil {
		for k, v := range cookies {
			request.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}

	response, err := a.Client.Do(request)
	if err != nil {
		return nil, err
	}
	if a.DebugMode {
		log.Printf("request=%v\n", *request)
		log.Printf("response=%v\n", *response)
	}
	return response, ParseResponse(a, response.Body, responseStruct)
}
