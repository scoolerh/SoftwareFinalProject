package game

import (
	"fmt"
	"strings"
)

func Testhandler() {
	fmt.Println("success")
}

func Steve(possibleMoves []MoveType, color string) MoveType {
	if color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)-1]
	}
}

func Joe(possibleMoves []MoveType, playerColor string, g Game) MoveType {
	var opponentColor string
	var bestMove MoveType = possibleMoves[0]
	var bestScore float32 = 0
	if playerColor == "w" {
		opponentColor = "b"
	} else {
		opponentColor = "w"
	}

	for index, move := range possibleMoves {
		_ = index
		originalState := g.State
		tempCaptured := make(map[string]int)
		tempCaptured["w"] = g.Captured["w"]
		tempCaptured["b"] = g.Captured["b"]
		var score float32

		tempState := g.UpdateState(playerColor, move)
		g.State = originalState

		for i := 1; i <= 24; i++ {
			slot := tempState[i]
			if strings.Contains(slot, "w") && strings.Contains(slot, "b") {
				move.CapturePiece = true
			}
		}

		//update the temporary captured map
		if move.CapturePiece {
			tempCaptured[opponentColor] += 1
		}

		pips := countPips(tempState, tempCaptured)
		score = 0.01 * float32(pips[opponentColor]-pips[playerColor])

		blots, _ := countBlots(tempState)
		score += float32(blots[opponentColor] - blots[playerColor])

		towers, _ := countTowers(tempState)
		score += float32(towers[playerColor] - towers[opponentColor])

		//chooses the first move to achieve the best score. Prioritizes moving pieces from the back.
		if playerColor == "w" {
			if score > bestScore {
				bestScore = score
				bestMove = move
			}
		} else {
			if score >= bestScore {
				bestScore = score
				bestMove = move
			}
		}
	}
	return bestMove
}

// can remove all the indicies stuff
func countBlots(gameState [26]string) (map[string]int, []int) {
	blots := make(map[string]int)
	var blotsIndices []int
	for index, slot := range gameState {
		if len(slot) == 1 {
			blots[slot] += 1
			blotsIndices = append(blotsIndices, index)
		} else if len(slot) > 1 && slot[0:1] != slot[1:2] {
			blots[slot[1:2]] += 1 //accounting for blots happening when piece is captured
		}
	}
	return blots, blotsIndices
}

func countTowers(gameState [26]string) (map[string]int, []int) {
	towers := make(map[string]int)
	var towersIndices []int
	for index, slot := range gameState {
		if len(slot) > 1 && slot[0:1] == slot[1:2] {
			towers[slot[0:1]] += 1
			towersIndices = append(towersIndices, index)
		}
	}
	return towers, towersIndices
}
