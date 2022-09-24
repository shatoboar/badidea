package app

import (
	"testing"
)

func setupUsers(t *testing.T) map[int]*User {
	t.Helper()
	users := make(map[int]*User, 0)
	for i := 0; i < 10; i++ {
		user := &User{
			UserId: i,
			Score:  i,
		}
		users[i] = user
	}
	return users
}
func TestGetTopUser(t *testing.T) {
	users := setupUsers(t)

	length := 10
	top := getTopUsers(users, length)
	for i := 0; i < length-1; i++ {
		t.Logf("test")
		if users[top[i]].Score < users[top[i+1]].Score {
			t.Fatalf("Expected user %d to have a better score than user %d", top[i+1], top[i])
		}
	}

}

func TestGetTopUserWIthDuplicate(t *testing.T) {
	users := setupUsers(t)

	users[5].Score = 8
	users[1].Score = 8

	length := 10
	top := getTopUsers(users, length)
	for i := 0; i < length-1; i++ {
		t.Logf("test")
		if users[top[i]].Score < users[top[i+1]].Score {
			t.Fatalf("Expected user %d to have a better score than user %d", top[i+1], top[i])
		}
	}

}
func TestGetTopUserWIthOnlyDuplicate(t *testing.T) {
	users := setupUsers(t)

	users[0].Score = 8
	users[1].Score = 8
	users[2].Score = 8
	users[3].Score = 8
	users[4].Score = 8
	users[5].Score = 8
	users[6].Score = 8
	users[7].Score = 8
	users[8].Score = 8
	users[9].Score = 8

	length := 10
	top := getTopUsers(users, length)
	for i := 0; i < length-1; i++ {
		t.Logf("test")
		if users[top[i]].Score < users[top[i+1]].Score {
			t.Fatalf("Expected user %d to have a better score than user %d", top[i+1], top[i])
		}
	}

}
