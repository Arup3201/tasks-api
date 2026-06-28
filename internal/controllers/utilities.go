package controllers

import (
	"errors"
	"net/http"
)

var (
	ErrMissingUserID = errors.New("user_id is missing in the request")
)

func getUserID(req *http.Request) (string, error) {
	id, ok := req.Context().Value("user_id").(string)
	if !ok {
		return "", ErrMissingUserID
	}

	return id, nil
}
