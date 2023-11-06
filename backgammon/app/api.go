package main

import (
	"backgammon/game"
	"fmt"
	"html/template"
	"net/http"
)

var games []game.Game
var initialState = [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}
var testState = [26]string{"", "ww", "bb", "w", "b", "ww", "bb", "w", "b", "ww", "bb", "w", "", "", "", "", "b", "ww", "bb", "w", "b", "ww", "bb", "w", "b", ""}
var p1 game.Player
var p2 game.Player
var whoseTurn string = "first"
var gameid int
var winner string

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Print the rules and how to use the tool for the user
func help(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "./html/index.html")
}

// todo: Create a database for users, allow a user to log in (or sign up if they do not have a username)
func login(writer http.ResponseWriter, req *http.Request) {
	//username := req.URL.Query().Get("user")
	http.ServeFile(writer, req, "./html/login.html")
}

// Starts a new game for the user and displays the initial board
func newgame(writer http.ResponseWriter, req *http.Request) {
	p1, p2 = game.Player{Id: "STEVE", Color: "w"}, game.Player{Id: "JOE", Color: "b"} //will need to be an input in the future
	gameid = len(games)
	capturedMap := initializeCapturedMap()
	g := game.Game{Gameid: gameid, Player1: p1, Player2: p2, State: initialState, Captured: capturedMap}
	games = append(games, g)
	variables := map[string]interface{}{"id": gameid, "p1": p1.Id, "p2": p2.Id}
	outputHTML(writer, "./html/newgame.html", variables)
}

func initializeCapturedMap() map[string]int {
	m := make(map[string]int)
	m["w"] = 0
	m["b"] = 0
	return m
}

func testplay(writer http.ResponseWriter, req *http.Request) {

	//for testing purposes
	p1, p2 := game.Player{Id: "STEVE", Color: "w"}, game.Player{Id: "JOE", Color: "b"} //will need to be an input in the future
	gameid := len(games)
	capturedMap := initializeCapturedMap()
	g := game.Game{Gameid: gameid, Player1: p1, Player2: p2, State: testState, Captured: capturedMap}
	games = append(games, g)
	fmt.Fprint(writer, "TIME TO PLAY \n")

	fmt.Fprintf(writer, "%v \n", g.State)
	for i := 0; i < 100; i++ {
		if g.IsWon() != "" {
			fmt.Fprint(writer, "WINNER")
			return
		}

		/* // returning and printing boardState for testing purposes
		log.Printf("\n move nr %v: \n", i)
		fmt.Fprintf(writer, "move nr %v: \n", i)
		g.Move(g.Player1)
		//fmt.Fprintf(writer, "player 1 made a move: %v", g.currMove)
		fmt.Fprintf(writer, "%v", g.State)
		fmt.Fprintf(writer, "captured pieces: %v \n", g.Captured)
		g.Move(g.Player2)
		//fmt.Fprintf(writer, "player 2 made a move: %v", g.currMove)
		fmt.Fprintf(writer, "%v", g.State)
		fmt.Fprintf(writer, "captured pieces: %v \n", g.Captured) */
	}
}

// Check whose turn it is and if the game is won, then have the player make a move
func play(writer http.ResponseWriter, req *http.Request) {
	g := games[gameid]
	if g.IsWon() != "" {
		if whoseTurn == "w" {
			winner = p2.Id
		}
		if whoseTurn == "b" {
			winner = p1.Id
		}
		won(writer, req)
		return
	}
	if whoseTurn == "first" {
		variables := map[string]interface{}{"id": gameid, "player": p1.Id, "captured": g.Captured, "one": g.State[0], "two": g.State[1], "three": g.State[2], "four": g.State[3], "five": g.State[4], "six": g.State[5], "seven": g.State[6], "eight": g.State[7], "nine": g.State[8], "ten": g.State[9], "eleven": g.State[10], "twelve": g.State[11], "thirteen": g.State[12], "fourteen": g.State[13], "fifteen": g.State[14], "sixteen": g.State[15], "seventeen": g.State[16], "eighteen": g.State[17], "nineteen": g.State[18], "twenty": g.State[19], "twentyone": g.State[20], "twentytwo": g.State[21], "twentythree": g.State[22], "twentyfour": g.State[23], "twentyfive": g.State[24], "twentysix": g.State[25]}
		outputHTML(writer, "./html/playing.html", variables)
		whoseTurn = "w"
		games[gameid] = g
	} else if whoseTurn == "w" {
		variables := map[string]interface{}{"id": gameid, "player": p2.Id, "captured": g.Captured, "one": g.State[0], "two": g.State[1], "three": g.State[2], "four": g.State[3], "five": g.State[4], "six": g.State[5], "seven": g.State[6], "eight": g.State[7], "nine": g.State[8], "ten": g.State[9], "eleven": g.State[10], "twelve": g.State[11], "thirteen": g.State[12], "fourteen": g.State[13], "fifteen": g.State[14], "sixteen": g.State[15], "seventeen": g.State[16], "eighteen": g.State[17], "nineteen": g.State[18], "twenty": g.State[19], "twentyone": g.State[20], "twentytwo": g.State[21], "twentythree": g.State[22], "twentyfour": g.State[23], "twentyfive": g.State[24], "twentysix": g.State[25]}
		g.Move(g.Player1)
		outputHTML(writer, "./html/playing.html", variables)
		whoseTurn = "b"
		games[gameid] = g
	} else {
		variables := map[string]interface{}{"id": gameid, "player": p1.Id, "captured": g.Captured, "one": g.State[0], "two": g.State[1], "three": g.State[2], "four": g.State[3], "five": g.State[4], "six": g.State[5], "seven": g.State[6], "eight": g.State[7], "nine": g.State[8], "ten": g.State[9], "eleven": g.State[10], "twelve": g.State[11], "thirteen": g.State[12], "fourteen": g.State[13], "fifteen": g.State[14], "sixteen": g.State[15], "seventeen": g.State[16], "eighteen": g.State[17], "nineteen": g.State[18], "twenty": g.State[19], "twentyone": g.State[20], "twentytwo": g.State[21], "twentythree": g.State[22], "twentyfour": g.State[23], "twentyfive": g.State[24], "twentysix": g.State[25]}
		g.Move(g.Player2)
		outputHTML(writer, "./html/playing.html", variables)
		whoseTurn = "w"
		games[gameid] = g
	}
}

// todo: if someone has won, update the database with wins/losses for each player. Print final board.
func won(writer http.ResponseWriter, req *http.Request) {
	variables := map[string]interface{}{"winner": winner}
	outputHTML(writer, "./html/won.html", variables)
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
	http.HandleFunc("/testplay", testplay)
	http.HandleFunc("/login", login)
	http.HandleFunc("/won", won)
	http.HandleFunc("/scoreboard", scoreboard)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
