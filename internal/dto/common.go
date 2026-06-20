package dto

import "fmt"

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }

func errField(msg string) error { return &ValidationError{Msg: msg} }

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(msg string, data interface{}) Response {
	return Response{Message: msg, Data: data}
}

func Fail(msg string) Response {
	return Response{Message: msg}
}

type PaginationQuery struct {
	Page  int
	Limit int
}

func (p *PaginationQuery) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Limit
}

func (p *PaginationQuery) SafeLimit() int {
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 10
	}
	return p.Limit
}

// suppress unused warning
var _ = fmt.Sprintf
