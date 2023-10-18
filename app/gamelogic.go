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
	var moveToDo MoveType
	var possibleMoves []MoveType
	currState := g.State

	if currPlayer == "w" {
		for i := 1; i <= 24; i++ {
			if strings.Contains(currState[i], "w") {
				for index, die := range dice {
					if 25-i >= die {
						goalPlace := currState[i+die]
						if !(strings.Contains(goalPlace, "b") && len(goalPlace) >= 2) {
							moveToDo.Slot = i
							moveToDo.Die = die
							moveToDo.DieIndex = index
							moveToDo.CapturePiece = false
							possibleMoves = append(possibleMoves, moveToDo)
						}
					}
				}
			}
		}
	} else if currPlayer == "b" {
		for i := 1; i <= 24; i++ {

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

//func capturePiece (slot, color)
// have some kind of object to keep track of captured pieces. Could be struct captured with w: and b: or just a list/slice of w's and b's
//if checkForCapturedPiece is true, add piece to captured pieces, and remove it from the board (so updateboard)
//right now updatestate takes in a move. Either change updatestate to take in pice, from, too, or make a separate update function,
//or make a move that
//reflects capturing, and make a way to handle that in updateState

//TODO:
//create capturePiece
//create checkForCapturedPiece
//update the updateState function
//call capturePiece in the appropriate place - decide where you want to call checkForCapturedPiece

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
	State    [26]string
	Captured map[string]int //so like w:2, b:3. Make sure to make this whenever we intialize Game!
	//only have one type player, in getmove have an if-statement that checks for human or AI, then execute different versions
	//NOTE that this is currently wrong. We need this to be a player, but either human or AI, i dont know how to do that
	//needs to not be an ai, but a player, a general human or ai - i think we might need a player struct...
	//state map[string]string //maps a string to an int, kind of like dictionary in python. Could also use array for this.

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
