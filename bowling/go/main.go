package bowling

import (
	"strconv"
	"strings"
)

func pointsForGame(game string) int {
	arr := gameToSlice(game)
	score := 0
	rounds := 1
	for len(arr) > 0 {
		if rounds >= 10 {
			if checkStrike(arr[0]) {
				score += 10
				arr = arr[1:]
			} else if len(arr) >= 2 && checkSpare(arr[0], arr[1]) {
				score += 10
				arr = arr[2:]
			} else {
				score += pointsForThrow(arr[0])
				arr = arr[1:]
			}
			continue
		}
		if checkStrike(arr[0]) {
			score += handleStrike(arr)
			arr = arr[1:]
		} else if checkSpare(arr[0], arr[1]) { //Check for spare
			score += handleSpare(arr)
			arr = arr[2:]
		} else {
			score += pointsForRound(arr[0], arr[1])
			arr = arr[2:]
		}
		rounds += 1
	}
	return score
}

func gameToSlice(game string) []string {
	return strings.Split(game, "")
}

//Check for strikes, takes a single args
func checkStrike(throw string) bool {
	if throw == "X" {
		return true
	}
	return false
}

func checkSpare(t1 string, t2 string) bool {
	if t2 == "/" {
		return true
	}
	return false
}

//Count the points gotten for a non strike round, which has two throws
func pointsForRound(t1 string, t2 string) int {
	if t2 == "/" {
		return 10
	}
	return pointsForThrow(t1) + pointsForThrow(t2)
}

//Count points for a single throw which is not a strike or a spare
func pointsForThrow(t string) int {
	if t == "-" {
		return 0
	} else if t == "X" {
		return 10
	}
	n, _ := strconv.Atoi(t)
	return n
}

//Handle spare scenario until a game with no spare and strike is made
func handleSpare(arr []string) int {
	sum := 10

	//Check if there is a 3rd throw and add the points for the third throw
	if len(arr) >= 3 {
		sum += pointsForThrow(arr[2])
	}

	return sum
}

//Count the score for the strike round
func handleStrike(arr []string) int {
	sum := 10

	// Check next 2 rounds to see if they are strikes
	for i := 1; i < 3; i++ {

		//Make sure that the index is valid
		if i > len(arr)-1 {
			break
		}

		if checkStrike(arr[i]) {
			sum += 10
			continue
		}

		//Spare only works if its the next round
		if i == 1 && checkSpare(arr[i], arr[i+1]) {
			sum += 10
		} else {
			sum += pointsForThrow(arr[i])
		}

		break
	}

	return sum
}
