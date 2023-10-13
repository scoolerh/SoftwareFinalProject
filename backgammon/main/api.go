package main

import (
	"backgammon/gamelogic"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var games []gamelogic.Game //will be a valid type when we fix packages
var initialState = [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}

// Print the rules and how to use the tool for the user
func help(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprint(writer, "Welcome to BACKGAMMON. <br>To start a new game, use the /newgame/ route. This will display the starting board and your game id.<br>To see how many wins or losses you have, you can use the /scoreboard/(player)/ route, where player is your username. For Matt, your username is Matt.<br>Enjoy!!!")
}

// todo: Create a database for users, allow a user to log in (or sign up if they do not have a username)
func login(writer http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("user")
	fmt.Fprint(writer, "Welcome "+username+" !")
}

// Starts a new game for the user and displays the initial board
func newgame(writer http.ResponseWriter, req *http.Request) {
	p1, p2 := 0, 0 //will need to be an input in the future
	gameid := len(games)
	game := gamelogic.Game(gameid, p1, p2, initialState) //will be fixed when packages are fixed
	games = append(games, game)
	fmt.Fprint(writer, initialState) //needs to be sent to HTML template
}

// todo: Check whose turn it is, if the game is won, have the player make a move
func play(writer http.ResponseWriter, req *http.Request) {
	//make sure to get the gameid through the url
	gameStr := req.URL.Query().Get("gameid")
	gameid, errGameid := strconv.Atoi(gameStr)
	if errGameid != nil {
		log.Fatal(errGameid) //might want to change the way to handle errors
	}
	game := games[gameid]

	//plays 10 moves
	for i := 0; i < 10; i++ {
		// returning and printing boardState for testing purposes
		boardState := game.move(game.player1)
		fmt.Fprint(writer, "player 1 made a move: "+boardState)
		boardState = game.move(game.player2)
		fmt.Fprint(writer, "player 2 made a move: "+boardState)
	}
	//does this need to return something?
}

//HOW THE REST OF THIS API LOOKS DEPENDS ON STRUCTURE, HTML ETC
//BRYAN AND HANNAH, CAN YOU TAKE RESPONSIBILITY FOR THIS?

// for i := 0; i < 10; i++ {
// 	move()
// }
// for i in range(10):
//     #returning and printing boardState for testing purposes
//     boardState = game.move(game.player1)
//     print("player 1 made a move: ", boardState)
//     boardState = game.move(game.player2)
//     print("player 2 made a move: ", boardState)
// currentState = json.dumps(game.board.currentState)
// return requests.get(f"http://frontend:5000/displayboard/{currentState}/").text

// todo: if someone has won, update the database with wins/losses for each player. Print final board.
func won(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprint(writer, "Hannah won!")
}

// todo: set up SQL database, check if the user is an actual user in the db, then return their win/loss ratio.”'
func scoreboard(writer http.ResponseWriter, req *http.Request) {
	//display player's win/loss ratio
	user := "Hannah"
	request := "http://db:5000/getprofile/" + user //currently in Flask
	fmt.Fprint(writer, request)
	//fmt.Fprint(writer, requests.get(request).text)
}

func main() {
	http.HandleFunc("/", help) //this makes an endpoint that calls the help function
	http.HandleFunc("/newgame", newgame)
	http.HandleFunc("/play", play)
	http.HandleFunc("/login", login)
	http.HandleFunc("/won", won)
	http.HandleFunc("/scoreboard", scoreboard)
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 5555, with standard mapping
}
