package client

import "fmt"

type DocmostError struct {
	StatusCode int
	Message    string
}

func (e *DocmostError) Error() string {
	return fmt.Sprintf("docmost: HTTP %d: %s", e.StatusCode, e.Message)
}
