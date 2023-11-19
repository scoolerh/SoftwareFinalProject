package game

import (
	"math/rand"
	"strings"
)

type Game struct {
	Gameid   string
	Player1  Player
	Player2  Player
	CurrTurn Player
	State    [26]string
	Captured map[string]int
	Pips     map[string]int
	Dice     []int
}

type MoveType struct {
	Slot, Die, DieIndex int
	CapturePiece        bool
}

type Player struct {
	Id    string
	Color string
}

// gets all the possible moves the player can choose from
func (g Game) GetPossibleMoves(dice []int, currPlayer string) []MoveType {
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
							if canBearOff || goalSlot != 25 {
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

		if len(possibleMoves) == 0 && canBearOff {
			for index, die := range dice {
				for i := 25 - die; i < 25; i++ {
					if strings.Contains(currState[i], "w") {
						move.Slot = i
						move.Die = 25 - i
						move.DieIndex = index
						move.CapturePiece = false
						possibleMoves = append(possibleMoves, move)
						//this should be a forced move, only one possibility, so return right away
						return possibleMoves
					}
				}
			}
		}

		//same process for black.
		//Note that black moves in opposite direction of white, so bearing of slot, home board and direction of dice are all different
	} else if currPlayer == "b" {
		canBearOff := g.isBearingOffAllowed("b")
		if g.Captured["b"] == 0 {
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
		} else {
			for index, die := range dice {
				goalState := currState[25-die]
				if !(strings.Contains(goalState, "w") && len(goalState) >= 2) {
					move.Slot = 25
					move.Die = -die
					move.DieIndex = index
					move.CapturePiece = false
					possibleMoves = append(possibleMoves, move)
				}
			}
		}
		if len(possibleMoves) == 0 && canBearOff {
			for index, die := range dice {
				for i := die; i > 0; i-- {
					if strings.Contains(currState[i], "b") {
						move.Slot = i
						move.Die = -i
						move.DieIndex = index
						move.CapturePiece = false
						possibleMoves = append(possibleMoves, move)
						//this should be a forced move, only one possibility, so return right away
						return possibleMoves
					}
				}
			}
		}
	}
	possibleMoves = removeDuplicateMoves(possibleMoves)
	return possibleMoves
}

// not a perfect function for moving duplicates, but works for this use: 
//either all the dice are the same, or they are all unique
func removeDuplicateMoves(possibleMoves []MoveType) []MoveType {

	checkedmap := make(map[int]int)
	var newMoves []MoveType

	for _, move := range possibleMoves {
		val, status := checkedmap[move.Slot]
		if !status || val != move.Die {
			checkedmap[move.Slot] = move.Die
			newMoves = append(newMoves, move)
		}
	}

	return newMoves
}

// returns empty string if nobody has won. Else, returns "b" or "w"
func (g Game) IsWon() string {
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

// checks if any pieces are not in the player's home section of the board
func (g Game) isBearingOffAllowed(playerColor string) bool {
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

func (g *Game) UpdateCaptured(move MoveType) {
	if g.CurrTurn.Color == "w" && move.Slot == 0 {
		g.Captured["w"] -= 1
	} else if g.CurrTurn.Color == "b" && move.Slot == 25 {
		g.Captured["b"] -= 1
	}
}

// deletes the die that was just played
func (g *Game) UpdateDice(dieIndex int) {
	g.Dice = DeleteElement(g.Dice, dieIndex)
}

// updates the state of the board to reflect most recent move
func (g *Game) UpdateState(playerColor string, move MoveType) [26]string {
	currState := g.State
	die := move.Die

	originalSpace := move.Slot
	newSpace := originalSpace + die
	originalSpaceState := currState[originalSpace]

	//removing piece from original space
	if originalSpace != 0 && originalSpace != 25 { //pieces with these locations are either captured or beared off
		currState[originalSpace] = originalSpaceState[0 : len(originalSpaceState)-1]
	}

	//if piece in endSlot is captured there will only be one piece there
	if move.CapturePiece {
		currState[newSpace] = playerColor
	} else {
		currState[newSpace] = currState[newSpace] + playerColor
	}

	g.State = currState
	return g.State
}

// switches turns when necesesary
func (g *Game) UpdateTurn() {
	if g.CurrTurn == g.Player1 {
		g.CurrTurn = g.Player2
	} else {
		g.CurrTurn = g.Player1
	}
}

// creates the game
// the player set as currturn here will play second, not first
func CreateGame(games []Game, user1 string, user2 string) (Game, [26]string) {
	p1, p2 := Player{Id: user1, Color: "w"}, Player{Id: user2, Color: "b"}
	initialState := [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}
	capturedMap := initializeCapturedMap()
	game := Game{Player1: p1, Player2: p2, CurrTurn: p2, State: initialState, Captured: capturedMap}
	return game, initialState
}

func countPips(gameState [26]string, capturedPieces map[string]int) map[string]int {
	pips := make(map[string]int)
	//goes through board to add pips for all pieces
	for i := 1; i <= 24; i++ {
		slot := gameState[i]
		if strings.Contains(slot, "w") {
			//each piece of a slot must move i spots
			pips["w"] += len(slot) * i
		} else if strings.Contains(slot, "b") {
			pips["b"] += len(slot) * (25 - i)
		}
	}
	//each captured piece must move 25 pips
	pips["w"] += capturedPieces["w"] * 25
	pips["b"] += capturedPieces["b"] * 25

	return pips
}

// from tutorialspoint.com
// deletes an element in a slice
func DeleteElement(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

// creating a map that will contain information on each player's number of captured pieces
func initializeCapturedMap() map[string]int {
	m := make(map[string]int)
	m["w"] = 0
	m["b"] = 0
	return m
}

// gets the move from the AI, either steve or joe
func GetAIMove(possibleMoves []MoveType, player Player, game Game) MoveType {
	var move MoveType
	if player.Id == "steve" { //AI 1
		move = Steve(possibleMoves, player.Color)
	} else if player.Id == "joe" { //AI 2
		move = Joe(possibleMoves, player.Color, game)
	}
	return move
}

// helper function to roll dice
func RollDice(numDice int) []int {
	var dice []int
	for i := 0; i < numDice; i++ {
		die := rand.Intn(6) + 1
		dice = append(dice, die)
	}
	if dice[0] == dice[1] {
		dice = append(dice, dice...)
	}
	return dice
}

// checks if there is a piece that is captured if move is made
func WillCapturePiece(endSlotState string, playerColor string) bool {
	return len(endSlotState) == 1 && endSlotState != playerColor
	//so if the length of state is 1 and the color is not the same as the moving piece, it is captured
}
