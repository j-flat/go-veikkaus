package goveikkaus

import (
	"fmt"
	"net/http"

	api "go-veikkaus/internal/veikkausapi"
)

type SessionCookies struct {
	JSessionID string `json:"JSESSIONID"`
}

type LoginPayload struct {
	Type     string `json:"type"`
	User     string `json:"login"`
	Password string `json:"password"`
}

type LoginSuccessful = map[string]interface{}

type LoginService service

func (s *LoginService) Login(username, password string) (LoginSuccessful, error) {
	client := http.Client{}
	payloadStruct := LoginPayload{
		Type:     "STANDARD_LOGIN",
		User:     username,
		Password: password,
	}

	body, err := api.GetJsonPayload(payloadStruct)

	if err != nil {
		fmt.Println("Could not get response payload:", err)
		return nil, err
	}

	req, err := api.GetRequest(api.LoginEndpoint, api.Post, body)

	if err != nil {
		fmt.Println("Could not get request object to login to the service. Error was:", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error login in to service. ERR:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if !api.ResponseCodeIsOk(resp.StatusCode) {
		return nil, fmt.Errorf("API returned a non-successful response. Status Code: %d", resp.StatusCode)
	}

	var result = LoginSuccessful{}
	err = api.HandleResponse(resp, result)

	if err != nil {
		fmt.Println("Oh no!")
	}

	if len(result) == 0 {
		fmt.Println("Login successful")
	} else {
		fmt.Println("Login unsuccesful:", result)
	}

	return nil, nil
}
