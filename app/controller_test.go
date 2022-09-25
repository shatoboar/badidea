package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

type TestVerifier struct {
}

func (tv *TestVerifier) verifyUser(r *http.Request) bool {
	return true
}

func setup(t *testing.T) *Server {
	t.Helper()

	s := &Server{
		DB:     NewDB(),
		Router: mux.NewRouter(),
		Auth:   &TestVerifier{},
	}
	s.RegisterRoutes()
	return s
}

var testUser = &User{
	UserId:        12345,
	UserName:      "dannyG",
	PickupHistory: []*Trash{},
	ReportHistory: []*Trash{},
	Rank:          2,
	Score:         0,
}

var testTrash = &Trash{
	ID:           [16]byte{},
	Latitude:     52.520008,
	Longitude:    13.404954,
	ImageURL:     "",
	ReportedBy:   0,
	ReportNumber: 0,
	Reward:       1,
}

func TestCreateUser(t *testing.T) {
	s := setup(t)
	recorder := httptest.NewRecorder()

	body, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("failed to encode the body: %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewReader(body))
	s.CreateUser(recorder, req)

	if recorder.Result().StatusCode != http.StatusCreated {
		t.Fatalf("Expected %d, got %d", http.StatusOK, recorder.Result().StatusCode)
	}
}

func TestNewTrash(t *testing.T) {
	s := setup(t)

	s.DB.Users[testUser.UserId] = testUser

	recorder := httptest.NewRecorder()
	body, err := json.Marshal(testTrash)
	if err != nil {
		t.Fatalf("Failed to encode the trash: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/trash", bytes.NewReader(body))
	req.Header.Add("user_id", strconv.Itoa(testUser.UserId))
	s.CreateNewTrash(recorder, req)

	gotStatus := recorder.Result().StatusCode

	if gotStatus != http.StatusCreated {
		t.Fatalf("Expected %d, got %d", http.StatusCreated, gotStatus)
	}
}

func TestReportTrash(t *testing.T) {
	s := setup(t)
	s.DB.Users[testUser.UserId] = testUser

	recorder := httptest.NewRecorder()
	body, err := json.Marshal(testTrash)
	if err != nil {
		t.Fatalf("Failed to encode the trash: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/trash", bytes.NewReader(body))
	req.Header.Set("user_id", strconv.Itoa(testUser.UserId))
	s.ReportTrash(recorder, req)

	got := recorder.Result().StatusCode
	if got != http.StatusCreated {
		t.Fatalf("Expected %d as status, got %d", http.StatusCreated, got)
	}

	for _, val := range s.DB.Trash {
		if val.Latitude == val.Longitude {
			t.Fatalf("Expected %v, got %v", testTrash, val)
		}
	}

	for _, val := range s.DB.Users {
		if val.Score != testTrash.Reward {
			t.Fatalf("Expected user %s to have %d score, but user has %d", testUser.UserName, testTrash.Reward, val.Score)
		}
	}
}
