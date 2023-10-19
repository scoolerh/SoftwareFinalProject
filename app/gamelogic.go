package main

import (
	"log"
	"math/rand"
	"strings"
)

type Gameplay interface {
	Move() map[string]string      //move should call getMove and doMove
	GetPossibleMoves() []MoveType //returns an array of moves (slot, die)
	UpdateState()                 //updates the state of the game to reflect the recent move
	IsWon() bool
	isBearingOffAllowed(playerColor string) bool //make capital if needed outside of package
}

func (g Game) IsWon() string {
	// returns empty string if nobody has won. Else, returns "b" or "w"
	gamestate := g.State
	//checks if all pieces are beared off
	if len(gamestate[0]) == 15 {
		return "b"
	} else if len(gamestate[25]) == 15 {
		return "w"
	} else {
		return ""
	}
}

func (g Game) isBearingOffAllowed(playerColor string) bool {
	//checks if any pieces are not in home board
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
	//administers everything that is needed to identify and execute move.
	//Changes done here lately might have to be reflected in changes to updateBoard
	log.Println("rolling dice")
	dice := RollDice(2) //change to input?
	for i := 0; i < len(dice); i++ {
		log.Println("getting possible move for this set of dice")
		possibleMoves := g.GetPossibleMoves(dice, player.Color)
		if len(possibleMoves) == 0 {
			log.Println("No possible moves, returning empty")
			return
		}
		log.Println("letting player choose move")
		move := GetMove(possibleMoves, player)
		if player.Color == "b" {
			dice[move.DieIndex] = -move.Die
		} else if player.Color == "w" {
			dice[move.DieIndex] = move.Die
		}

		log.Println("Checks if piece needs to be captured")
		//may want to abstract better later - fix when everything else is working (this might need to be redone when play endpoint is done)
		endSlot, endSlotState := getEndSlot(move, g.State) //check whether we want this to return both slot and state or just one
		_ = endSlot
		if willCapturePiece(endSlotState) {
			move.CapturePiece = true
			g.Captured["player.Color"] += 1
		}

		log.Println("updates gamestate according to move")
		g.UpdateState(player.Color, move)
		dice = DeleteElement(dice, move.DieIndex)
	}
}

func (g Game) UpdateState(playerColor string, move MoveType) [26]string {
	//updates the state of the board to reflect most recent move
	currState := g.State
	die := move.Die
	//call getEndSpace where that is applicable
	originalSpace := move.Slot
	newSpace := originalSpace + die
	originalSpaceState := currState[originalSpace] //change this variable name? //check indexing +/- 1 error

	if move.CapturePiece {
		currState[newSpace] = playerColor
	} else {
		// currState[newSpace] = originalSpaceState[0 : len(originalSpaceState)+1] //check indexing
		// newSpaceState := currState[newSpace]
		// currState[newSpace] = newSpaceState + playerColor

		// the three above lines could probably be simplified to
		currState[newSpace] = currState[newSpace] + playerColor
	}

	g.State = currState
	return g.State
}

func (g Game) GetPossibleMoves(dice []int, currPlayer string) []MoveType {
	//gets all the possible moves the player can choose from
	var move MoveType
	var possibleMoves []MoveType
	currState := g.State

	//logic for finding white's moves
	if currPlayer == "w" {
		//checks if the player is allowed to bear off pieces. Uses this information later.
		canBearOff := g.isBearingOffAllowed("w")
		//finds all possible moves when there are no captured pieces
		if g.Captured["w"] == 0 {
			//loops through all slots of board
			for i := 1; i <= 24; i++ {
				//locates slots where white has pieces
				if strings.Contains(currState[i], "w") {
					for index, die := range dice {
						//checks that the move won't move the piece off the board
						if 25-i >= die {
							goalSlot := i + die
							goalState := currState[i+die]
							//checks that either bearing off is legal, or that we are not planning on bearing off
							if canBearOff || goalSlot != 0 {
								//checks that the goal slot is not occupied by tower of opposite color
								if !(strings.Contains(goalState, "b") && len(goalState) >= 2) {
									//gets necessary numbers and adds move to list
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
			//if we have captured pieces, those are the only ones that can move
		} else {
			for index, die := range dice {
				//captured pieces start from 0
				goalState := currState[die]
				//does not need to check for bearing off, it is not possible when starting from 0
				if !(strings.Contains(goalState, "b") && len(goalState) >= 2) {
					move.Slot = 0
					move.Die = die
					move.DieIndex = index
					move.CapturePiece = false
					possibleMoves = append(possibleMoves, move)
				}
			}
		}

		//same process for black.
		//Note that black moves in opposite direction of white, so bearing of slot, home board and direction of dice are all different
	} else if currPlayer == "b" {
		canBearOff := g.isBearingOffAllowed("b")
		if g.Captured["b"] == 0 {
			for i := 1; i <= 24; i++ {
				for i := 1; i <= 24; i++ {
					if strings.Contains(currState[i], "b") {
						for index, die := range dice {
							if i >= die {
								goalSlot := i - die
								goalState := currState[i-die]
								if canBearOff || goalSlot != 0 {
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
	//helper function to get the end slot and end state of a move (where the piece is moving to)
	//is not used a lot yet. Should be used more later when we improve levels of abstraction
	originalSlot := move.Slot
	dieRoll := move.Die
	endSlot := originalSlot + dieRoll
	endSlotState := gameState[endSlot]
	return endSlot, endSlotState
}

func willCapturePiece(endSlotState string) bool {
	//checks if there is a piece thats captured if move is made
	return len(endSlotState) == 1
}

func RollDice(numDice int) []int {
	//helper function to roll dice
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
	//prompts either the player or the AI to pick a move
	var move MoveType
	if player.Id == 0 { //AI
		move = GetAIMove(possibleMoves, player.Color)
	} else { //human
		move = GetHumanMove(possibleMoves, player.Color)
	}

	return move
}

func GetHumanMove(possibleMoves []MoveType, color string) MoveType {
	//dummy function
	if color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)]
	}
}

// change this implementation
func GetAIMove(possibleMoves []MoveType, color string) MoveType {
	//picks the first possible move to do. Will be improved in the future
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
