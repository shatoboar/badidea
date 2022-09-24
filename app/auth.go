package app

import (
	"errors"
	"net/http"
	"strconv"
)

func verifyUser(r *http.Request) bool {
	return true
}

func decodeUserID(r *http.Request) (int, error) {
	user := r.Header.Get("user_id")
	if len(user) == 0 {
		return 0, errors.New("Expected header value, got none.")
	}
	userID, err := strconv.Atoi(user)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
