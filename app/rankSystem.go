package app

var Ranks [5]string = [...]string{"Rookie Hunter", "Hunter", "Commited Hunter", "Veteran Hunter", "True Trash Hunter"}
var RankMaxScore = map[string]int{
	"Rookie Hunter":     100,
	"Hunter":            500,
	"Commited Hunter":   1500,
	"Veteran Hunter":    5000,
	"True Trash Hunter": 10000,
}

var nextRank = map[string]string{
	"Rookie Hunter":     "Hunter",
	"Hunter":            "Commited Hunter",
	"Commited Hunter":   "Veteran Hunter",
	"Veteran Hunter":    "True Trash Hunter",
	"True Trash Hunter": "True Trash Hunter",
}

//increments score and update rank if uprank is possible
func updateRank(user *User, points int) {
	if user.Rank.Score > RankMaxScore[user.Rank.Title] {
		user.Rank.Title = nextRank[user.Rank.Title]
	}

	user.Rank.Score += points

}

//updates the leaderboard after updating a users score
func updateLeaderboard() {

}
