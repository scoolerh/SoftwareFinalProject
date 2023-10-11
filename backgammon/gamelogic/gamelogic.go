package gamelogic //does this need to be in its own directory or something?
//TODO: fix package structure

import (
	"math/rand"
	"strings"
)

/*
Questions for Matt:
	- How to handle player - human - ai without inheritance, things that need a player that can be either human or AI
	- How does the rest of the go structure look?
	- Is our thin thread really a thin thread?
		- We will have to change the structure to get moves from the player
		- Can we get some help with how to think about this?
*/

type gameplay interface {
	//all these still need to be implemented
	move() map[string]string //move should call getMove and doMove
	//rollDice() []int          //returns an array of ints. I dont think this has to be a part of the gameplay interface
	getPossibleMoves() []moveType //returns an array of moves (slot, die)
	updateState()                 //updates the state of the game to reflect the recent move
}

func (g Game) move(player player) {
	dice := rollDice(2) //change to input?
	for i := 0; i < len(dice); i++ {
		possibleMoves := getPossibleMoves(g, dice, player)
		if len(possibleMoves) == 0 {
			return g.state
		}
		move := getMove(possibleMoves, player)
		if player.color == "b" {
			die = -move[1]
		} else if player.color == "w" {
			die = move[1]
		}
		currState := doMove()
		dice := deleteElement(dice, move[2])
	}
}

func (g Game) updateState(playerColor string, move moveType) [26]string {
	currState := g.state
	die := move.die
	//figure out how to do two moves
	originalSpace := move.slot
	newSpace := originalSpace + die
	originalSpaceState := currState[originalSpace] //change this variable name? //check indexing +/- 1 error
	currState[newSpace] = originalSpaceState[0:-1] //can we do -1
	newSpaceState := currState[newSpace]
	currState[newSpace] = newSpaceState + playerColor
	g.state = currState
	return g.state
}

func (g Game) getPossibleMoves(dice []int, currPlayer string) []move {
	var moveToDo moveType
	var possibleMoves []moveType
	currState := g.state

	if currPlayer == "w" {
		for i := 1; i <= 24; i++ {
			if strings.Contains(currState[string(i)], "w") { //needs to use int to string converter here, this creates runes
				for index, die := range dice { //we don't need index, but it is returned. How to handle this?
					if 25-i >= die {
						goalPlace := currState[string(i+die)]
						if !(strings.Contains(goalPlace, "b") && len(goalPlace) >= 2) {
							moveToDo.slot = i
							moveToDo.die = die
							moveToDo.dieIndex = index
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

func rollDice(numDice int) []int {
	var dice []int
	for i := 0; i < numDice; i++ {
		die := rand.Intn(6) + 1
		dice = append(dice, die)
	}
	return dice
}

type Game struct {
	gameid           int
	player1, player2 player //only have one type player, in getmove have an if-statement that checks for human or AI, then execute different versions
	//NOTE that this is currently wrong. We need this to be a player, but either human or AI, i dont know how to do that
	//needs to not be an ai, but a player, a general human or ai - i think we might need a player struct...
	//state map[string]string //maps a string to an int, kind of like dictionary in python. Could also use array for this.
	state [26]string
}

// type gamestate struct { //could use this or a map for the gamestate - example of map is at the bottom
// 	tile0, tile1, tile2, tile3, tile4, tile5, tile6, tile7, tile8, tile9, tile10, tile11, tile12, tile13, tile14, tile15, tile16, tile17,
// 	tile18, tile19, tile20, tile21, tile22, tile23, tile24, tile25 string //this might need to be improved...
// }

type moveType struct {
	slot, die, dieIndex int
}

// // the following three will have to change
// type player interface {
// 	getMove() move
// }

type player struct {
	id    int //check if it best if this is string or int. Note that we might need to use a stringToInt fcn to convert
	color string
}

// type ai struct {
// 	id, color string
// }

func getMove(possibleMoves []moveType, player player) moveType {
	var move moveType
	if player.id == 0 { //AI
		move = getAIMove(possibleMoves, player.color)
	} else { //human
		move = getHumanMove(possibleMoves, player.color)
	}

	return move
}

func getHumanMove(possibleMoves []moveType, color string) moveType {
	if color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)]
	}
}

// change this implementation
func getAIMove(possibleMoves []moveType, color string) moveType {
	if color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)]
	}
}

/* 	// Create a map with string keys and int values
myMap := make(map[string]int)

// Assign values to keys
myMap["one"] = 1
myMap["two"] = 2
myMap["three"] = 3

// Access values using keys
fmt.Println("Value for key 'two':", myMap["two"])

// Check if a key exists
value, exists := myMap["four"]
if exists {
	fmt.Println("Value for key 'four':", value)
} else {
	fmt.Println("Key 'four' does not exist.")
}

initialState := map[string]string{"0": "", "1": "ww", "2": "", "3": "", "4": "", "5": "", "6": "bbbbbb", "7": "", "8": "bbb", "9": "", "10": "", "11": "", "12": "wwwwww",
	// "13": "bbbbb", "14": "", "15": "", "16": "", "17": "www", "18": "", "19": "wwwww", "20": "", "21": "", "22": "", "23": "", "24": "bb", "25": ""}
*/
