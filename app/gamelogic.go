package main

import (
	"math/rand"
	"strings"
)

type Gameplay interface {
	//all these still need to be implemented
	Move() map[string]string      //move should call getMove and doMove
	GetPossibleMoves() []MoveType //returns an array of moves (slot, die)
	UpdateState()                 //updates the state of the game to reflect the recent move
	IsWon() bool
	isMovingHomePossible(playerColor string) bool //make capital if needed outside of package
}

func (g Game) IsWon() string {
	// returns empty string if nobody has won. Else, returns "b" or "w"
	gamestate := g.State
	if len(gamestate[0]) == 15 {
		return "b"
	} else if len(gamestate[25]) == 15 {
		return "w"
	} else {
		return ""
	}
}

func (g Game) isMovingHomePossible(playerColor string) bool {
	gameState := g.State
	var outPieces int
	if playerColor == "w" {
		for i := 1; i <= 18; i++ {
			if strings.Contains(gameState[i], "w") {
				outPieces += 1
			}
		}
		if outPieces == 0 && g.Captured["w"] == 0 {
			return true
		} else {
			return false
		}
	} else {
		for i := 7; i <= 24; i++ {
			if strings.Contains(gameState[i], "b") {
				outPieces += 1
			}
		}
		if outPieces == 0 && g.Captured["b"] == 0 {
			return true
		} else {
			return false
		}
	}
}

// currently returning nothing. originally returned game state but i don't think we need to
func (g Game) Move(player Player) {
	dice := RollDice(2) //change to input?
	for i := 0; i < len(dice); i++ {
		possibleMoves := g.GetPossibleMoves(dice, player.Color)
		if len(possibleMoves) == 0 {
			return
		}
		move := GetMove(possibleMoves, player)
		if player.Color == "b" {
			dice[move.DieIndex] = -move.Die
		} else if player.Color == "w" {
			dice[move.DieIndex] = move.Die
		}

		//may want to abstract better later - fix when everything else is working
		endSlot, endSlotState := getEndSlot(move, g.State) //check whether we want this to return both slot and state or just one
		_ = endSlot
		if checkForCapturedPiece(endSlotState) {
			move.CapturePiece = true
			g.Captured["player.Color"] += 1
		}

		g.UpdateState(player.Color, move)
		dice = DeleteElement(dice, move.DieIndex)
	}
}

func (g Game) UpdateState(playerColor string, move MoveType) [26]string {
	currState := g.State
	die := move.Die
	//call getEndSpace where that is applicable
	originalSpace := move.Slot
	newSpace := originalSpace + die
	originalSpaceState := currState[originalSpace] //change this variable name? //check indexing +/- 1 error

	if move.CapturePiece {
		currState[newSpace] = playerColor
	} else {
		currState[newSpace] = originalSpaceState[0 : len(originalSpaceState)+1] //check indexing
		newSpaceState := currState[newSpace]
		currState[newSpace] = newSpaceState + playerColor

		// the three above lines could probably be simplified to
		//currState[newSpace] = currState[newSpace] + playerColor
	}

	g.State = currState
	return g.State
}

func (g Game) GetPossibleMoves(dice []int, currPlayer string) []MoveType {
	//run as normally, then remove illegal moves if not isMovingHomePossible? Or check first?
	var move MoveType
	var possibleMoves []MoveType
	currState := g.State

	if currPlayer == "w" {
		canMovePiecesHome := g.isMovingHomePossible("w")
		if g.Captured["w"] == 0 {
			for i := 1; i <= 24; i++ {
				if strings.Contains(currState[i], "w") {
					for index, die := range dice {
						if 25-i >= die {
							goalSlot := i + die
							goalState := currState[i+die]
							if canMovePiecesHome || goalSlot != 0 {
								if !(strings.Contains(goalState, "b") && len(goalState) >= 2) {
									move.Slot = i
									move.Die = die
									move.DieIndex = index
									move.CapturePiece = false
									possibleMoves = append(possibleMoves, move)
								}
							}
						}
					}
				}
			}
		} else {
			for index, die := range dice {
				goalState := currState[die]
				if !(strings.Contains(goalState, "b") && len(goalState) >= 2) {
					move.Slot = 0
					move.Die = die
					move.DieIndex = index
					move.CapturePiece = false
					possibleMoves = append(possibleMoves, move)
				}
			}
		}
	} else if currPlayer == "b" {
		canMovePiecesHome := g.isMovingHomePossible("b")
		if g.Captured["b"] == 0 {
			for i := 1; i <= 24; i++ {
				for i := 1; i <= 24; i++ {
					if strings.Contains(currState[i], "b") {
						for index, die := range dice {
							if i >= die {
								goalSlot := i - die
								goalState := currState[i-die]
								if canMovePiecesHome || goalSlot != 0 {
									if !(strings.Contains(goalState, "w") && len(goalState) >= 2) {
										move.Slot = i
										move.Die = -die
										move.DieIndex = index
										move.CapturePiece = false
										possibleMoves = append(possibleMoves, move)
									}
								}
							}
						}
					}
				}
			}
		} else {
			for index, die := range dice {
				goalState := currState[die]
				if !(strings.Contains(goalState, "w") && len(goalState) >= 2) {
					move.Slot = 25
					move.Die = -die
					move.DieIndex = index
					move.CapturePiece = false
					possibleMoves = append(possibleMoves, move)
				}
			}
		}
	}

	return possibleMoves
}

func getEndSlot(move MoveType, gameState [26]string) (int, string) {
	originalSlot := move.Slot
	dieRoll := move.Die
	endSlot := originalSlot + dieRoll
	endSlotState := gameState[endSlot]
	return endSlot, endSlotState
}

func checkForCapturedPiece(endSlotState string) bool {
	return len(endSlotState) == 1
}

func RollDice(numDice int) []int {
	var dice []int
	for i := 0; i < numDice; i++ {
		die := rand.Intn(6) + 1
		dice = append(dice, die)
	}
	return dice
}

type Game struct {
	Gameid   int
	Player1  Player
	Player2  Player
	CurrTurn Player
	State    [26]string
	Captured map[string]int
}

type MoveType struct {
	Slot, Die, DieIndex int
	CapturePiece        bool
}

type Player struct {
	Id    int //check if it best if this is string or int. Note that we might need to use a stringToInt fcn to convert
	Color string
}

func GetMove(possibleMoves []MoveType, player Player) MoveType {
	var move MoveType
	if player.Id == 0 { //AI
		move = GetAIMove(possibleMoves, player.Color)
	} else { //human
		move = GetHumanMove(possibleMoves, player.Color)
	}

	return move
}

func GetHumanMove(possibleMoves []MoveType, color string) MoveType {
	if color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)]
	}
}

// change this implementation
func GetAIMove(possibleMoves []MoveType, color string) MoveType {
	if color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)]
	}
}

// from tutorialspoint.com
func DeleteElement(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}
