package gateway

type Result struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type ResultError struct {
	Status int   `json:"status"`
	Error  error `json:"error"`
}
