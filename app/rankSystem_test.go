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

func TestUpdateRank(t *testing.T) {
	users := setupUsers(t)

	users[0].Score = 100
	users[0].Title = "Rookie Hunter"

	users[1].Score = 500
	users[1].Title = "Hunter"

	users[2].Score = 50000
	users[2].Title = "True Trash Hunter"

	updateRank(users[0], 1)
	updateRank(users[1], 10)
	updateRank(users[2], 10)

	if users[0].Score != 101 || users[0].Title != "Hunter" {
		t.Fatalf("Expected Score %d to have value '101' and Title %s to be equal to 'Hunter'", users[0].Score, users[0].Title)
	}

	if users[1].Score != 510 || users[1].Title != "Committed Hunter" {
		t.Fatalf("Expected Score %d to have value '510' and Title %s to be equal to 'Committed Hunter'", users[1].Score, users[1].Title)
	}

	if users[2].Score != 50010 || users[2].Title != "True Trash Hunter" {
		t.Fatalf("Expected Score %d to have value '50010' and Title %s to be equal to 'True Trash Hunter'", users[2].Score, users[2].Title)
	}
}
