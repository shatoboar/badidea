package app

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"firebase.google.com/go/auth"
)

type AuthClient struct {
	Client *auth.Client
}

func (ac *AuthClient) verifyUser(r *http.Request) bool {
	token := r.Header.Get("jwt_token")
	_, err := ac.Client.VerifyIDToken(context.Background(), token)
	if err != nil {
		return false
	}

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
