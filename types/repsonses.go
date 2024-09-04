package types

import "fmt"

type GenericResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Err error
}

func (s *GenericResponse) Response() string {
	return fmt.Sprintf("Message: %v", s.Message)
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v:", e.Err)
}
