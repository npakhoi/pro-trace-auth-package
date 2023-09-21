package refresh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IRefreshFunction interface {
	ResetToken(rfToken string) Response
}

type function struct {
	host string
}

func NewRefreshService(host string) IRefreshFunction {
	return function{host: host}
}

func (f function) ResetToken(rfToken string) Response {
	url := f.host + "/refresh-token"
	rfTokenO := TokenRefreshRq{
		RefreshToken: rfToken,
	}
	requestBody, err := json.Marshal(rfTokenO)
	if err != nil {
		return Response{
			Data:   TokenRefreshRs{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error marshaling JSON: %v", err),
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return Response{
			Data:   TokenRefreshRs{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error creating HTTP request: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{
			Data:   TokenRefreshRs{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error sending HTTP request: %v", err),
		}
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			Data:   TokenRefreshRs{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error reading response body: %v", err),
		}
	}

	var res Response

	if err := json.Unmarshal(body, &res); err != nil {
		return Response{
			Data:   TokenRefreshRs{},
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("Error parsing JSON: %v", err),
		}
	}
	return res
}
