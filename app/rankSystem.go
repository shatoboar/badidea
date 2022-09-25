package app

import "sort"

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

var Ranks [5]string = [...]string{"Rookie Hunter", "Hunter", "Commited Hunter", "Veteran Hunter", "True Trash Hunter"}
var RankMaxScore = map[string]int{
	"Rookie Hunter":    100,
	"Hunter":           500,
	"Committed Hunter": 1500,
	"Veteran Hunter":   5000,
}

var nextRank = map[string]string{
	"Rookie Hunter":    "Hunter",
	"Hunter":           "Committed Hunter",
	"Committed Hunter": "Veteran Hunter",
	"Veteran Hunter":   "True Trash Hunter",
}

func contains(s []int, id int) bool {
	for _, v := range s {
		if v == id {
			return true
		}
	}

	return false
}

// increments score and update rank if uprank is possible
func updateRank(user *User, points int) {
	user.Score += points
	if user.Score > RankMaxScore[user.Title] && user.Title != "True Trash Hunter" {
		user.Title = nextRank[user.Title]
	}
}

func getTopUsers(users map[string]*User, top int) []*User {
	topUsers := make([]*User, 0)

	keys := make([]string, 0, len(users))

	for k := range users {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return users[keys[i]].Score < users[keys[j]].Score
	})

	for _, k := range keys {
		topUsers = append(topUsers, users[k])
		if len(topUsers) == top {
			break
		}
	}

	return topUsers
}
