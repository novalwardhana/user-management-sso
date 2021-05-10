package model

type Result struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type Authorize struct {
	UniqueCode string `json:"unique_code"`
}
