package main

import (
	"fmt"
	"net/http"
	"text/template"
	// "encoding/json" // Maybe don't need this
)

var games []Game //will be a valid type when we fix packages
var initialState = [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}

// Print the rules and how to use the tool for the user
func help(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "Frontend/html/index.html")
}

// todo: Create a database for users, allow a user to log in (or sign up if they do not have a username)
func login(writer http.ResponseWriter, req *http.Request) {
	username := req.URL.Query().Get("user")
	fmt.Fprint(writer, "Welcome "+username+" !")
}

// Creates a new game for the user
// TODO: Create a button that sends user to /play endpoint with corresponding game id
// TODO: Get user input about players, initial state
// TODO: Get new player IDs with each new game
func newgame(writer http.ResponseWriter, req *http.Request) {
	p1, p2 := Player{Id: 1, Color: "w"}, Player{Id: 2, Color: "b"} //will need to be an input in the future
	gameid := len(games)
	game := Game{Gameid: gameid, Player1: p1, Player2: p2, CurrTurn: p1, State: initialState}
	games = append(games, game)
	http.ServeFile(writer, req, "Frontend/html/game.html")
}

// Assume correct game-id is given.
func play(writer http.ResponseWriter, req *http.Request) {
	possibleMoves := game.GetPossibleMoves(dice, game.CurrTurn.Color)
	tmpl := template.Must(template.ParseFiles("Frontend/html/play.html"))
	err := tmpl.Execute(writer, possibleMoves)
	if err != nil {
		panic(err)
	}
	http.ServeFile(writer, req, "Frontend/html/play.html")
	//make sure to get the gameid through the url
	// gameStr := req.URL.Query().Get("gameid")
	// gameid, errGameid := strconv.Atoi(gameStr)
	// if errGameid != nil {
	// 	log.Fatal(errGameid) //might want to change the way to handle errors
	// }
	// game := games[gameid]

	// // var gameWon bool 	//default is False

	// // do 10 moves
	// //for loop placeholder for testing. should be: for !gameWon {
	// for i := 0; i < 10; i++ {
	// 	dice := RollDice(2)
	// 	for i := 0; i < len(dice); i++ {
	// 		tmpl := template.Must(template.ParseFiles("Frontend/html/play.html"))
	// 		possibleMoves := game.GetPossibleMoves(dice, game.CurrTurn.Color)
	// 		if len(possibleMoves) == 0 {
	// 			fmt.Fprint(writer, "no possible moves")
	// 			break //tell player there are no moves, go to next turn
	// 		}
	// 		var move MoveType
	// 		if game.CurrTurn.Id != 0 {
	// 			err := tmpl.Execute(writer, possibleMoves)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			http.ServeFile(writer, req, "Frontend/html/play.html")
	// 			reqBody, err := ioutil.ReadAll(req.Body)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			fmt.Println(reqBody) // for testing

	// 			// move := GetHumanMove(possibleMoves, game.CurrTurn.Color) //TODO: get user input
	// 		} else {
	// 			http.ServeFile(writer, req, "Frontend/html/play.html")
	// 			// move := GetAIMove(possibleMoves, game.CurrTurn.Color)
	// 		}

	// 		if game.CurrTurn.Color == "b" {
	// 			dice[i] = -move.Die //was originally move.DieIndex
	// 		} else if game.CurrTurn.Color == "w" {
	// 			dice[i] = move.Die //was originally move.DieIndex
	// 		}
	// 		game.UpdateState(game.CurrTurn.Color, move)
	// 		dice = DeleteElement(dice, move.DieIndex)
	// 		http.ServeFile(writer, req, "Frontend/html/play.html")
	// 	}
	// 	if game.CurrTurn == game.Player1 {
	// 		game.CurrTurn = game.Player2
	// 	} else {
	// 		game.CurrTurn = game.Player1
	// 	}
	// 	//gameWon := IsWon()
	// }
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
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 5555, with standard mapping
}
