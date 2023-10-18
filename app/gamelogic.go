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

		g.UpdateState(player.Color, move)
		dice = DeleteElement(dice, move.DieIndex)
	}
}

func (g Game) UpdateState(playerColor string, move MoveType) [26]string {
	currState := g.State
	die := move.Die
	//figure out how to do two moves
	originalSpace := move.Slot
	newSpace := originalSpace + die
	originalSpaceState := currState[originalSpace]                          //change this variable name? //check indexing +/- 1 error
	currState[newSpace] = originalSpaceState[0 : len(originalSpaceState)+1] //check indexing
	newSpaceState := currState[newSpace]
	currState[newSpace] = newSpaceState + playerColor
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

func checkForCapturedPiece(move MoveType) bool {
	//if the slot we are moving to/have moved to (be very aware of this!) has one of the different color, return true
	//not sure if this needs to be separate, or just part of capturePiece
	return false
}

//func capturePiece (slot, color)
// have some kind of object to keep track of captured pieces. Could be struct captured with w: and b: or just a list/slice of w's and b's
//if checkForCapturedPiece is true, add piece to captured pieces, and remove it from the board (so updateboard)
//right now updatestate takes in a move. Either change updatestate to take in pice, from, too, or make a separate update function, or make a move that
//reflects capturing, and make a way to handle that in updateState

func RollDice(numDice int) []int {
	var dice []int
	for i := 0; i < numDice; i++ {
		die := rand.Intn(6) + 1
		dice = append(dice, die)
	}
	return dice
}

type Game struct {
	Gameid  int
	Player1 Player
	Player2 Player
	State   [26]string
	//only have one type player, in getmove have an if-statement that checks for human or AI, then execute different versions
	//NOTE that this is currently wrong. We need this to be a player, but either human or AI, i dont know how to do that
	//needs to not be an ai, but a player, a general human or ai - i think we might need a player struct...
	//state map[string]string //maps a string to an int, kind of like dictionary in python. Could also use array for this.

}

type MoveType struct {
	Slot, Die, DieIndex int
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
