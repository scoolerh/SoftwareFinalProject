package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

// WE NEED TO HANDLE "EMPTY" URL KEYWORDS

// var games []Game //will be a valid type when we fix packages
// var games = append(games, createTestGame(1, p1))
var initialState = [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}
var testState = [26]string{"", "ww", "bb", "w", "b", "ww", "bb", "w", "b", "ww", "bb", "w", "", "", "", "", "", "b", "ww", "bb", "w", "b", "ww", "bb", "w", "b"}
var p1 Player
var p2 Player
var whoseTurn string = "first"
var gameid int
var winner string

// var games = []Game{createTestGame(0, p1)}
var game = createTestGame(0, p1)

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
	// p1, p2 = Player{Id: "STEVE", Color: "w"}, Player{Id: "JOE", Color: "b"} //will need to be an input in the future
	// gameid = len(games)
	// capturedMap := initializeCapturedMap()
	// var dice []int
	// game := Game{Gameid: gameid, Player1: p1, Player2: p2, State: initialState, Captured: capturedMap, Dice: dice}
	// games = append(games, game)
	urlParams := url.Values{}
	strValues := ConvertParams(50, 0, 0, false)
	strGameid := strconv.Itoa(game.Gameid)
	urlParams.Add("gameid", strGameid)
	urlParams.Add("Slot", strValues[0])
	startGameURL := "/play?" + urlParams.Encode()
	variables := map[string]interface{}{"id": game.Gameid, "p1": p1.Id, "p2": p2.Id, "startGameURL": startGameURL}
	outputHTML(writer, "./html/newgame.html", variables)
}

// func initializeCapturedMap() map[string]int {
// 	m := make(map[string]int)
// 	m["w"] = 0
// 	m["b"] = 0
// 	return m
// }

// func testplay(writer http.ResponseWriter, req *http.Request) {

// 	//for testing purposes
// 	p1, p2 := Player{Id: "STEVE", Color: "w"}, Player{Id: "JOE", Color: "b"} //will need to be an input in the future
// 	gameid := len(games)
// 	capturedMap := initializeCapturedMap()
// 	game := Game{Gameid: gameid, Player1: p1, Player2: p2, State: testState, Captured: capturedMap}
// 	games = append(games, game)
// 	fmt.Fprint(writer, "TIME TO PLAY \n")

// 	fmt.Fprintf(writer, "%v \n", game.State)
// 	for i := 0; i < 10; i++ {
// 		// returning and printing boardState for testing purposes
// 		log.Printf("\n move nr %v: \n", i)
// 		fmt.Fprintf(writer, "move nr %v: \n", i)
// 		game.Move(game.Player1)
// 		fmt.Fprintf(writer, "player 1 made a move: %v", game.currMove)
// 		fmt.Fprintf(writer, "%v \n", game.State)
// 		//fmt.Fprintf(writer, "captured pieces: %v \n", game.Captured)
// 		game.Move(game.Player2)
// 		fmt.Fprintf(writer, "player 2 made a move: %v", game.currMove)
// 		fmt.Fprintf(writer, "%v \n", game.State)
// 		//fmt.Fprintf(writer, "captured pieces: %v \n", game.Captured)
// 	}
// }

// Check whose turn it is and if the game is won, then have the player make a move
// TODO: build URL for the template to send data to. Should be "/play?key=vlue" stuff
// TODO: implement so that the list of all games is accessed through database
// TODO: display dice
// TODO: tell user if no possible moves
func play(writer http.ResponseWriter, req *http.Request) {
	// u, err := url.Parse(s) //is this necessary
	// if err != nil {
	// 	panic(err)
	// }
	// urlVars, _ := url.ParseQuery(u.RawQuery) //parse query param into map
	urlVars := req.URL.Query()
	varGameid := urlVars["gameid"][0]
	gameid, _ := strconv.Atoi(varGameid)
	// g := games[gameid] // This needs to be changed to work with database
	g := game

	//if no move
	if urlVars["Slot"][0] != "50" {
		//abstract this?
		varSlot := urlVars["Slot"][0]
		varDie := urlVars["Die"][0]
		varDieIndex := urlVars["DieIndex"][0]
		varCapturePiece := urlVars["CapturePiece"][0]

		slot, _ := strconv.Atoi(varSlot)
		die, _ := strconv.Atoi(varDie)
		dieIndex, _ := strconv.Atoi(varDieIndex)
		capturePiece, _ := strconv.ParseBool(varCapturePiece)
		move := MoveType{Slot: slot,
			Die:          die,
			DieIndex:     dieIndex,
			CapturePiece: capturePiece,
		}
		g.updateDice(dieIndex)
		g.UpdateState(g.CurrTurn.Color, move)
		if g.IsWon() != "" {
			if whoseTurn == "w" {
				winner = p2.Id
			}
			if whoseTurn == "b" {
				winner = p1.Id
			}
			won(writer, req)
		}
	}

	//think about this logic for first turn
	g.updateTurn()

	//rolls the dice if the dice list is empty
	if len(g.Dice) == 0 {
		g.Dice = RollDice(2)
	}

	//display after dice roll
	//TODO: pass in dice?
	isHuman := true
	if g.CurrTurn.id == 0 {
		isHuman = false
	}
	var outputVars1 = map[string]interface{}{"id": gameid, "player": g.CurrTurn, "state": g.State, "captured": g.Captured, "isHuman": isHuman}
	outputHTML(writer, "./html/playing.html", outputVars1)

	possibleMoves := g.GetPossibleMoves(g.Dice, g.CurrTurn.Color)
	//deletes all dice if no possible moves
	//should somehow communicate to user that there are no possible moves
	if len(possibleMoves) == 0 {
		g.Dice = nil
	}
	var outputVars map[string]interface{}
	var human bool
	var playerId, _ = strconv.Atoi(g.CurrTurn.Id)
	if playerId != 0 {
		human = true
		var urlList []string
		if len(possibleMoves) == 0 {
			urlParams := url.Values{}
			strValues := ConvertParams(50, 0, 0, false)
			urlParams.Add("gameid", varGameid)
			urlParams.Add("Slot", strValues[0])
			var urlString string = "/play?" + urlParams.Encode()
			urlList = append(urlList, urlString)
		} else {
			for index, move := range possibleMoves {
				_ = index
				urlParams := url.Values{}
				strValues := ConvertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
				urlParams.Add("gameid", varGameid)
				urlParams.Add("Slot", strValues[0])
				urlParams.Add("Die", strValues[1])
				urlParams.Add("DieIndex", strValues[2])
				urlParams.Add("CapturePiece", strValues[3])
				var urlString string = "/play?" + urlParams.Encode()
				urlList = append(urlList, urlString)
			}

		}
		outputVars = map[string]interface{}{"possibleMoves": possibleMoves, "urlList": urlList, "game": g, "human": human}
	} else {
		human = false
		urlParams := url.Values{}
		urlParams.Add("gameid", varGameid)
		if len(possibleMoves) != 0 {
			move := GetAIMove(possibleMoves, g.CurrTurn.Color)
			strValues := ConvertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
			urlParams.Add("Slot", strValues[0])
			urlParams.Add("Die", strValues[1])
			urlParams.Add("DieIndex", strValues[2])
			urlParams.Add("CapturePiece", strValues[3])
		}
		url := "/play?" + urlParams.Encode()
		outputVars = map[string]interface{}{"url": url, "human": human}
	}
	outputHTML(writer, "./html/playing.html", outputVars)
}

// todo: if someone has won, update the database with wins/losses for each player. Print final board.
func won(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprint(writer, winner+" won!")
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
	// http.HandleFunc("/testplay", testplay)
	http.HandleFunc("/login", login)
	http.HandleFunc("/won", won)
	http.HandleFunc("/scoreboard", scoreboard)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
