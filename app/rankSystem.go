package app

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

var Ranks [5]string = [...]string{"Rookie Hunter", "Hunter", "Commited Hunter", "Veteran Hunter", "True Trash Hunter"}
var RankMaxScore = map[string]int{
	"Rookie Hunter":     100,
	"Hunter":            500,
	"Committed Hunter":  1500,
	"Veteran Hunter":    5000,
	"True Trash Hunter": 10000,
}

var nextRank = map[string]string{
	"Rookie Hunter":     "Hunter",
	"Hunter":            "Commited Hunter",
	"Committed Hunter":  "Veteran Hunter",
	"Veteran Hunter":    "True Trash Hunter",
	"True Trash Hunter": "True Trash Hunter",
}

func contains(s []int, id int) bool {
	for _, v := range s {
		if v == id {
			return true
		}
	}

	return false
}

//increments score and update rank if uprank is possible
func updateRank(user *User, points int) {
	user.Score += points
	if user.Score > RankMaxScore[user.Title] {
		user.Title = nextRank[user.Title]
	}
}

func getTopUsers(users map[int]*User, top int) []int {
	topUsers := make([]int, 0)
	maxScore, maxScoreKey := 0, 0

	for key, val := range users {
		if val.Score > maxScore {
			maxScore = val.Score
			maxScoreKey = key
		}
	}
	topUsers = append(topUsers, maxScoreKey)

	for i := 0; i < top-1; i++ {
		prevLargestValue := 0
		prevLargestKey := 0
		for key, val := range users {
			if val.Score >= prevLargestValue && val.Score <= maxScore && !contains(topUsers, key) {
				prevLargestValue = val.Score
				prevLargestKey = key

			}
		}
		topUsers = append(topUsers, prevLargestKey)
		maxScore = prevLargestValue
	}

	return topUsers
}
