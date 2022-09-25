package app

import (
	"context"
	"errors"
	"net/http"

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

func decodeUserName(r *http.Request) (string, error) {
	user := r.Header.Get("user_name")
	if len(user) == 0 {
		return "", errors.New("Expected header value, got none.")
	}
	return user, nil
}
