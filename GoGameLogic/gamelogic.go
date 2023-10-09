package gogamelogic //does this need to be in its own directory or something?

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
	move() gamestate //move should call getMove and doMove
	//rollDice() []int          //returns an array of ints. I dont think this has to be a part of the gameplay interface
	getPossibleMoves() []move //returns an array of moves (slot, die)
}

func (g game) getPossibleMoves(dice []int, currPlayer string) []move {
	var moveToDo move
	var possibleMoves []move
	currState := g.state

	if currPlayer == "w" {
		for i := 1; i <= 24; i++ {
			if strings.Contains(currState[string(i)], "w") {
				for index, die := range dice { //we don't need index, but it is returned. How to handle this?
					if 25-i >= die {
						goalPlace := currState[string(i+die)]
						if !(strings.Contains(goalPlace, "b") && len(goalPlace) >= 2) {
							moveToDo.slot = i
							moveToDo.die = die
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

type game struct {
	player1, player2 player
	//NOTE that this is currently wrong. We need this to be a player, but either human or AI, i dont know how to do that
	//needs to not be an ai, but a player, a general human or ai - i think we might need a player struct...
	state map[string]string //maps a string to an int, kind of like dictionary in python. Seems a little inconvenient, see example below
}

type gamestate struct { //could use this or a map for the gamestate - example of map is at the bottom
	tile0, tile1, tile2, tile3, tile4, tile5, tile6, tile7, tile8, tile9, tile10, tile11, tile12, tile13, tile14, tile15, tile16, tile17,
	tile18, tile19, tile20, tile21, tile22, tile23, tile24, tile25 string //this might need to be improved...
}

type move struct {
	slot, die int
}

type player interface {
	getMove() move
}

type human struct {
	id, color string
}

type ai struct {
	id, color string
}

func (player ai) getMove(possibleMoves []move) move {
	if player.color == "w" {
		return possibleMoves[0]
	} else {
		return possibleMoves[len(possibleMoves)]
	}
}

func (player human) getMove(possibleMoves []move) move {
	return possibleMoves[0] //change this implementation
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
*/
