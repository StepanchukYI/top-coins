package server

type Server interface {
	Start() error
}

type ErrorResponse struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}
