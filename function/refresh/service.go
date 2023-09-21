package refresh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IRefreshFunction interface {
	ResetToken(rfToken string) TokenRefreshRs
}

type function struct {
	host string
}

func NewRefreshService(host string) IRefreshFunction {
	return function{host: host}
}

func (f function) ResetToken(rfToken string) TokenRefreshRs {
	url := f.host + "/refresh-token"
	rfTokenO := TokenRefreshRq{
		RefreshToken: rfToken,
	}
	requestBody, err := json.Marshal(rfTokenO)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return TokenRefreshRs{}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return TokenRefreshRs{}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return TokenRefreshRs{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
		return TokenRefreshRs{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return TokenRefreshRs{}
	}

	var res Response

	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return TokenRefreshRs{}
	}
	return res.Data
}
