package gogamelogic //does this need to be in its own directory or something?

type gameplay interface {
	//all these still need to be implemented
	move() gamestate          //move should call getMove and doMove
	rollDice() []int          //returns an array of ints
	getPossibleMoves() []move //returns an array of moves (slot, die)
}

type game struct {
	player1, player2 player
	//NOTE that this is currently wrong. We need this to be a player, but either human or AI, i dont know how to do that
	state map[string]int //maps a string to an int, kind of like dictionary in python. Seems a little inconvenient, see example below
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

type gamestate struct { //could use this or a map for the gamestate
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
