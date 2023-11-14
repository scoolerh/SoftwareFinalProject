package game

import (
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

type Gameplay interface {
	Move() map[string]string      //move should call getMove and doMove
	GetPossibleMoves() []MoveType //returns an array of moves (slot, die)
	UpdateState()                 //updates the state of the game to reflect the recent move
	IsWon() bool
	isBearingOffAllowed(playerColor string) bool //make capital if needed outside of package
}

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
	Id    string //check if it best if this is string or int. Note that we might need to use a stringToInt fcn to convert
	Color string
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

func countPips(gameState [26]string, capturedPieces map[string]int) map[string]int { //can change this to just update pips based on move if necessary
	pips := make(map[string]int)
	//goes through board to add pips for all pieces
	for i := 1; i <= 24; i++ {
		slot := gameState[i]
		if strings.Contains(slot, "w") {
			//each piece of a slot must move i spots
			pips["w"] += len(slot) * i
		} else if strings.Contains(slot, "b") {
			pips["w"] += len(slot) * (25 - i)
		}
	}
	//each captured piece must move 25 pips
	pips["w"] += capturedPieces["w"] * 25
	pips["b"] += capturedPieces["b"] * 25

	return pips
}

func (g *Game) UpdateState(playerColor string, move MoveType) [26]string { //the * makes it a pointer and not a value. Remember this if similar issues arise later.
	//updates the state of the board to reflect most recent move
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

func (g Game) GetPossibleMoves(dice []int, currPlayer string) []MoveType {
	//gets all the possible moves the player can choose from
	var move MoveType
	var possibleMoves []MoveType
	currState := g.State

	//logic for finding white's moves
	if currPlayer == "w" {
		//checks if the player is allowed to bear off pieces. Uses this information later.
		canBearOff := g.isBearingOffAllowed("w")
		log.Printf("Bearing off status: %v", canBearOff)
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
		log.Printf("Bearing off status: %v", canBearOff)
		if g.Captured["b"] == 0 {
			for i := 1; i <= 24; i++ {
				if strings.Contains(currState[i], "b") {
					for index, die := range dice {
						if i >= die {
							goalSlot := i - die
							goalState := currState[i-die]
							if canBearOff || goalSlot != 0 { //could add prints for further ebugging
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

	return possibleMoves
}

func WillCapturePiece(endSlotState string, playerColor string) bool {
	//checks if there is a piece thats captured if move is made
	return len(endSlotState) == 1 && endSlotState != playerColor
	//so if the length of state is 1 and the color is not the same as the moving piece, it is captured
}

func RollDice(numDice int) []int {
	//helper function to roll dice
	var dice []int
	for i := 0; i < numDice; i++ {
		die := rand.Intn(6) + 1
		dice = append(dice, die)
	}
	if dice[0] == dice[1] { //assuming only two dice. Might be changed later if we want more
		dice = append(dice, dice...)
	}
	return dice
}

func GetMove(possibleMoves []MoveType, player Player, game Game) MoveType {
	//prompts either the player or the AI to pick a move
	var move MoveType
	if player.Id == "steve" { //AI
		move = Steve(possibleMoves, player.Color) //only one now, implement Joe later
	} else if player.Id == "joe" {
		move = Joe(possibleMoves, player.Color, game)
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
		return possibleMoves[len(possibleMoves)-1]
	}
}

// change this implementation
func GetAIMove(possibleMoves []MoveType, color string) MoveType {
	//picks the first possible move to do. Will be improved in the future
	if color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)-1]
	}
}

// from tutorialspoint.com
func DeleteElement(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

func ConvertParams(slot int, die int, index int, capture bool) [4]string {
	strSlot := strconv.Itoa(slot)
	strDie := strconv.Itoa(die)
	strDieIndex := strconv.Itoa(index)
	strCapturePiece := strconv.FormatBool(capture)
	returns := [4]string{strSlot, strDie, strDieIndex, strCapturePiece}
	return returns
}

func initializeCapturedMap() map[string]int {
	m := make(map[string]int)
	m["w"] = 0
	m["b"] = 0
	return m
}

// the player set as currturn here will play second, not first
func CreateGame(games []Game, user1 string, user2 string) (Game, [26]string) {
	p1, p2 := Player{Id: user1, Color: "w"}, Player{Id: user2, Color: "b"} //will need to be an input in the future
	initialState := [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}
	// testState := [26]string{"bbbbbbbbbbbbbb", "b", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "ww", "wwwwwwwwwwwww"}
	capturedMap := initializeCapturedMap()
	game := Game{Player1: p1, Player2: p2, CurrTurn: p2, State: initialState, Captured: capturedMap}
	return game, initialState
}

func ParseVariables(urlVariables url.Values) (int, int, int, bool) {
	varSlot := urlVariables["Slot"][0]
	log.Printf("the slot is: %v", varSlot)
	varDie := urlVariables["Die"][0]
	varDieIndex := urlVariables["DieIndex"][0]
	varCapturePiece := urlVariables["CapturePiece"][0]

	slot, _ := strconv.Atoi(varSlot)
	die, _ := strconv.Atoi(varDie)
	dieIndex, _ := strconv.Atoi(varDieIndex)
	capturePiece, _ := strconv.ParseBool(varCapturePiece)
	return slot, die, dieIndex, capturePiece
}

func AddUrlParams(urlParams url.Values, valuesToAdd [4]string) url.Values {
	urlParams.Add("Slot", valuesToAdd[0])
	urlParams.Add("Die", valuesToAdd[1])
	urlParams.Add("DieIndex", valuesToAdd[2])
	urlParams.Add("CapturePiece", valuesToAdd[3])
	return urlParams
}

// deletes the die that was just played
func (g *Game) UpdateDice(dieIndex int) {
	g.Dice = DeleteElement(g.Dice, dieIndex)
}

// switch turns if necesesary
func (g *Game) UpdateTurn() {
	if len(g.Dice) == 0 {
		if g.CurrTurn == g.Player1 {
			g.CurrTurn = g.Player2
		} else {
			g.CurrTurn = g.Player1
		}
	}
}

// currently returning nothing. originally returned game state but i don't think we need to
// func (g *Game) Move(player Player) {
// 	//administers everything that is needed to identify and execute move.
// 	//Changes done here lately might have to be reflected in changes to updateBoard
// 	dice := RollDice(2) //change to input?
// 	log.Printf("diceroll: %v \n", dice)
// 	numDice := len(dice)

// 	for i := 0; i < numDice; i++ {
// 		log.Printf("Using dice %v", i+1)

// 		possibleMoves := g.GetPossibleMoves(dice, player.Color)
// 		if len(possibleMoves) == 0 {
// 			log.Println("no possible moves")
// 			return
// 		}
// 		mockGame := *g
// 		move := GetMove(possibleMoves, player, mockGame)

// 		if player.Color == "b" {
// 			dice[move.DieIndex] = -move.Die
// 		} else if player.Color == "w" {
// 			dice[move.DieIndex] = move.Die
// 		}

// 		endSlot := move.Slot + move.Die
// 		endSlotState := g.State[endSlot]

// 		if willCapturePiece(endSlotState, player.Color) {
// 			move.CapturePiece = true
// 			g.Captured[endSlotState] += 1
// 		}
// 		log.Printf("player %s chose move %v", player.Color, move)
// 		//only for testing purposes
// 		g.currMove = move

// 		g.UpdateState(player.Color, move)
// 		log.Printf("state updated to: %v", g.State)
// 		g.Pips = countPips(g.State, g.Captured)

// 		if player.Color == "w" && move.Slot == 0 {
// 			g.Captured["w"] -= 1
// 		} else if player.Color == "b" && move.Slot == 25 {
// 			g.Captured["b"] -= 1
// 		}

// 		dice = DeleteElement(dice, move.DieIndex)
// 	}
// }
