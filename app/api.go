package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
)

// WE NEED TO HANDLE "EMPTY" URL KEYWORDS

var games []Game

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
	g := createGame(games)
	games = append(games, g)
	urlParams := url.Values{}
	strGameid := strconv.Itoa(g.Gameid)
	urlParams.Add("gameid", strGameid)
	urlParams.Add("Slot", "-1")
	startGameURL := "/play?" + urlParams.Encode()
	variables := map[string]interface{}{"id": g.Gameid, "p1": g.Player1.Id, "p2": g.Player2.Id, "startGameURL": startGameURL}
	outputHTML(writer, "./html/newgame.html", variables)
}

// Check whose turn it is and if the game is won, then have the player make a move
// TODO: display dice
func play(writer http.ResponseWriter, req *http.Request) {

	urlVars := req.URL.Query()
	varGameid := urlVars["gameid"][0]
	var outputVars map[string]interface{}
	var human bool
	// g := games[gameid] // This needs to be changed to work with database
	var noPossibleMoves bool
	var gameid, _ = strconv.Atoi(varGameid)
	g := games[gameid]

	//if there is a move
	if urlVars["Slot"][0] != "-1" {
		slot, die, dieIndex, capturePiece := parseVariables(urlVars)
		move := MoveType{Slot: slot,
			Die:          die,
			DieIndex:     dieIndex,
			CapturePiece: capturePiece,
		}

		endSlot := move.Slot + move.Die
		endSlotState := g.State[endSlot]
		if willCapturePiece(endSlotState, g.CurrTurn.Color) {
			move.CapturePiece = true
			g.Captured[endSlotState] += 1
		}
		fmt.Printf("player %s chose move %v \n", g.CurrTurn.Color, move)

		g.updateDice(dieIndex)
		if g.CurrTurn.Color == "w" && move.Slot == 0 {
			g.Captured["w"] -= 1
		} else if g.CurrTurn.Color == "b" && move.Slot == 25 {
			g.Captured["b"] -= 1
		}
		g.UpdateState(g.CurrTurn.Color, move)
		fmt.Printf("Board updated to: %v \n", g.State)
		fmt.Printf("dice left: %v \n", g.Dice)
		if g.IsWon() != "" {
			winner := g.CurrTurn.Id
			http.Redirect(writer, req, "/won?winner="+winner, http.StatusSeeOther)
		}
	} else {
		fmt.Println("no move")
	}

	//think about this logic for first turn
	g.updateTurn()
	var playerId, _ = strconv.Atoi(g.CurrTurn.Id)

	//rolls the dice if the dice list is empty
	if len(g.Dice) == 0 {
		g.Dice = RollDice(2)
		fmt.Printf("diceroll: %v \n", g.Dice)
	}

	possibleMoves := g.GetPossibleMoves(g.Dice, g.CurrTurn.Color)

	//deletes all dice if no possible moves
	if len(possibleMoves) == 0 {
		g.Dice = nil
		noPossibleMoves = true
	}

	if playerId != 0 {
		fmt.Println("human move now")
		human = true
		var urlList []string
		if len(possibleMoves) == 0 {
			urlParams := url.Values{}
			strValues := ConvertParams(-1, 0, 0, false)
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
				urlParams = addUrlParams(urlParams, strValues)
				var urlString string = "/play?" + urlParams.Encode()
				urlList = append(urlList, urlString)
			}
		}
		outputVars = map[string]interface{}{"possibleMoves": possibleMoves, "urlList": urlList, "game": g, "isHuman": human, "noPossibleMoves": noPossibleMoves, "state": g.State, "captured": g.Captured, "player": g.CurrTurn.Id}
	} else {
		fmt.Println("ai move now")
		human = false
		urlParams := url.Values{}
		urlParams.Add("gameid", varGameid)
		if len(possibleMoves) != 0 {
			move := GetAIMove(possibleMoves, g.CurrTurn.Color)
			strValues := ConvertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
			urlParams = addUrlParams(urlParams, strValues)
		} else {
			urlParams.Add("Slot", "-1")
		}
		url := "/play?" + urlParams.Encode()
		outputVars = map[string]interface{}{"url": url, "isHuman": human, "state": g.State, "captured": g.Captured, "player": g.CurrTurn.Id}
	}
	games[gameid] = g
	outputHTML(writer, "./html/playing.html", outputVars)
}

// todo: if someone has won, update the database with wins/losses for each player. Print final board.
// TODO: pass in winner
func won(writer http.ResponseWriter, req *http.Request) {
	winner := req.URL.Query().Get("winner")
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
	// http.HandleFunc("/testplay", testplay)
	http.HandleFunc("/login", login)
	http.HandleFunc("/won", won)
	http.HandleFunc("/scoreboard", scoreboard)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
