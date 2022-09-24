package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setup(t *testing.T) *Server {
	t.Helper()
	s := &Server{
		DB:     NewDB(),
		Router: mux.NewRouter(),
	}
	s.RegisterRoutes()
	return s
}

func TestCreateUser(t *testing.T) {
	s := setup(t)
	user := &User{
		UserId:        "test_id",
		UserName:      "dannyG",
		PickupHistory: []*Trash{},
		ReportHistory: []*Trash{},
		Rank:          "rookie",
		JWTToken:      "",
		FirebaseToken: "",
		Score:         1,
	}
	r := httptest.NewRecorder()

	body, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to encode the body: %v", err)
	}
	w := httptest.NewRequest(http.MethodGet, "/user", bytes.NewReader(body))
	s.CreateUser(r, w)

	if r.Result().StatusCode != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusOK, r.Result().StatusCode)
	}
}
