package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ILoginFunction interface {
	Login(cred Credentials) Response
}

type function struct {
	host string
}

func NewLoginService(host string) ILoginFunction {
	return function{host: host}
}

func (s function) Login(cred Credentials) Response {
	url := s.host + "/login"

	requestBody, err := json.Marshal(cred)
	if err != nil {
		return Response{
			Data:   LoginResponse{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error marshaling JSON: %v", err),
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return Response{
			Data:   LoginResponse{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error creating HTTP request: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{
			Data:   LoginResponse{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error sending HTTP request: %v", err),
		}
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			Data:   LoginResponse{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error reading response body: %v", err),
		}
	}

	var res Response

	if err = json.Unmarshal(body, &res); err != nil {
		return Response{
			Data:   LoginResponse{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error parsing JSON: %v", err),
		}
	}
	return res
}
