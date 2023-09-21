package common

type Response struct {
	Data   string `json:"data"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}
