package main

import (
	"fmt"
	"log"
	"net/http"
)

var games []Game //will be a valid type when we fix packages
var initialState = [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}

// Print the rules and how to use the tool for the user
func help(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "./html/index.html")
}

// todo: Create a database for users, allow a user to log in (or sign up if they do not have a username)
func login(writer http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("user")
	fmt.Fprint(writer, "Welcome "+username+" !")
}

// Starts a new game for the user and displays the initial board
func newgame(writer http.ResponseWriter, req *http.Request) {
	p1, p2 := Player{Id: 0, Color: "w"}, Player{Id: 0, Color: "b"} //will need to be an input in the future
	gameid := len(games)
	capturedMap := initializeCapturedMap()
	game := Game{Gameid: gameid, Player1: p1, Player2: p2, State: initialState, Captured: capturedMap}
	games = append(games, game)
	http.ServeFile(writer, req, "./html/game.html")
}

func initializeCapturedMap() map[string]int {
	m := make(map[string]int)
	m["w"] = 0
	m["b"] = 0
	return m
}

// todo: Check whose turn it is, if the game is won, have the player make a move
func play(writer http.ResponseWriter, req *http.Request) {
	//for testing purposes
	fmt.Fprint(writer, "TIME TO PLAY \n")

	game := games[0]
	fmt.Fprintf(writer, "%v \n", game.State)
	for i := 0; i < 50; i++ {
		// returning and printing boardState for testing purposes
		log.Println("calling move for player 1 and testing log")
		game.Move(game.Player1)
		//NOTE! GAME IS NOT PROPERLY UPDATED! See updateState
		fmt.Fprint(writer, "player 1 made a move")
		fmt.Fprintf(writer, "%v", game.State)
		fmt.Fprintf(writer, "captured pieces: %v \n", game.Captured)
		game.Move(game.Player2)
		fmt.Fprint(writer, "player 2 made a move")
		fmt.Fprintf(writer, "%v", game.State)
		fmt.Fprintf(writer, "captured pieces: %v \n", game.Captured)
	}

	/*original thing:
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
		game.Move(game.Player1)
		fmt.Fprint(writer, "player 1 made a move ")
		game.Move(game.Player2)
		//print the game state
		fmt.Fprint(writer, "player 2 made a move: ")
		//print the game state
	} */
}

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
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
