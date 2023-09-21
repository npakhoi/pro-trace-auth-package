package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ILoginFunction interface {
	Login(cred Credentials) LoginResponse
}

type function struct {
	host string
}

func NewLoginService(host string) ILoginFunction {
	return function{host: host}
}

func (s function) Login(cred Credentials) LoginResponse {
	url := s.host + "/login"

	requestBody, err := json.Marshal(cred)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return LoginResponse{}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return LoginResponse{}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return LoginResponse{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
		return LoginResponse{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return LoginResponse{}
	}

	var res Response

	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return LoginResponse{}
	}
	return res.Data
}
