package services

import (
	"fmt"
	"io"
	"net/http"
)

func (serv *Services) GetTextPlain(r *http.Request) (s string, err error) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("io.ReadAll(r.Body) fail: %w", err)
	}
	s = string(reqBody)
	return s, nil
}
