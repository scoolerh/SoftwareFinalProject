package game

import (
	"fmt"
	"log"
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

//need to be able to temporarily update state, then do countPips on that
//so countPips should be separate from a game fcn

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
		tempCaptured := g.Captured
		var score float32

		tempState := g.UpdateState(playerColor, move) //this does not execute cature pieces to not change to map
		g.State = originalState                       //need to do this because updatestate changes the game. Make sure this does not cause any issues.
		//Maybe we dont need it? We dont have *game - test this

		//check if the updatestate tried to capture piece
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

		pips := countPips(tempState, g.Captured)
		score = 0.01 * float32(pips[opponentColor]-pips[playerColor])
		log.Printf("pips: %v", pips)

		blots, blotsIndices := countBlots(tempState)
		score += float32(blots[opponentColor] - blots[playerColor]) //add some weight
		log.Printf("blots: %v", blots)
		log.Printf("blotindices: %v", blotsIndices)

		towers, towersIndices := countTowers(tempState)
		score += float32(towers[playerColor] - towers[opponentColor]) //add some weight
		log.Printf("towers: %v", towers)
		log.Printf("towerindices: %v", towersIndices)

		//could do checks for number of pieces in home board, do a strategy check to value leaving pieces behind etc, but this ai is fine for now

		if score > bestScore { //so we chose the first move to achieve the best score
			bestScore = score
			bestMove = move
			log.Printf("bestMove updated to %v", bestMove)
		}
	}
	return bestMove
}

func countBlots(gameState [26]string) (map[string]int, []int) {
	blots := make(map[string]int)
	var blotsIndices []int
	for index, slot := range gameState {
		//_ = index
		if len(slot) == 1 {
			blots[slot] += 1
			//log.Printf("blot found at %v", index)
			blotsIndices = append(blotsIndices, index)
		} else if len(slot) > 1 && slot[0:1] != slot[1:2] {
			blots[slot[1:2]] += 1 //you are doing this to account for blots happening when piece is captured. Do similar for tower too
		}
	}
	return blots, blotsIndices
}

func countTowers(gameState [26]string) (map[string]int, []int) {
	towers := make(map[string]int)
	var towersIndices []int
	for index, slot := range gameState {
		//_ = index
		if len(slot) > 1 && slot[0:1] == slot[1:2] {
			towers[slot[0:1]] += 1
			//log.Printf("tower found at %v", index)
			towersIndices = append(towersIndices, index)
		}
	}
	return towers, towersIndices
}
