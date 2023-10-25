package game

import (
	"fmt"
	"log"
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
	var bestMove MoveType
	var bestScore int = 0
	if playerColor == "w" {
		opponentColor = "b"
	} else {
		opponentColor = "w"
	}
	for index, move := range possibleMoves {
		_ = index
		originalState := g.State
		var score int
		tempState := g.UpdateState(playerColor, move)
		g.State = originalState //need to do this because updatestate changes the game. Make sure this does not cause any issues.
		//Maybe we dont need it? We dont have *game - test this
		pips := countPips(tempState, g.Captured)
		score = pips[opponentColor] - pips[playerColor] //we want high opponentpip and low playerpip

		//more checks to affect score here

		if score > bestScore { //so we chose the first move to achieve the best score
			bestScore = score
			bestMove = move
			log.Printf("bestMove updated")
		}
	}
	return bestMove
}
