package note

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type INoteFunction interface {
	NoteExpiredToken(token string)
}

type function struct {
	host string
}

func NewNoteFunction(host string) INoteFunction {
	return function{host: host}
}

func (f function) NoteExpiredToken(token string) {
	url := f.host + "/note-expired-token"
	rq := ExpiredTokenStruct{
		Token: token,
	}
	requestBody, err := json.Marshal(rq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
	}
}
