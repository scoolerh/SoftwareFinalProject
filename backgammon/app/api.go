package main

import (
	"backgammon/game"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "collective"
	dbname   = "backgammon"
)

var games []game.Game
var db *sql.DB
var currentUser string

//TODO: probably need a user variable here
//how does it look when two users are logged in? If two people play on different computers?
//and how does it look if two users play on the same computer? Is one or both going to log in?

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
	http.ServeFile(writer, req, "app/html/index.html")
}

// todo: Create a database for users, allow a user to log in (or sign up if they do not have a username)
func login(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("Connecting to login endpoint")
	http.ServeFile(writer, req, "app/html/login.html")
}

func register(writer http.ResponseWriter, req *http.Request) {
	var username string
	var password string

	if req.Method == http.MethodPost {
		username = req.FormValue("username")
		password = req.FormValue("password")
	}

	query := "INSERT INTO users VALUES ('" + username + "', '" + password + "')"

	var err error
	_, err = db.Exec(query)
	if err != nil {
		panic(err) //might want to change this later
	}
	http.ServeFile(writer, req, "app/html/index.html") //indicate somehow that registration was successful
}

func loggedin(writer http.ResponseWriter, req *http.Request) {
	var username string
	var password string

	if req.Method == http.MethodPost {
		username = req.FormValue("username")
		password = req.FormValue("password")
	}

	query := "SELECT password FROM users WHERE username='" + username + "'"

	rows, err := db.Query(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	var refPassword string
	for rows.Next() {
		rows.Scan(&refPassword)
	}

	if password != refPassword {
		http.ServeFile(writer, req, "app/html/loginfailed.html")
	}

	currentUser = username
	log.Printf("Welcome %s!", currentUser)
	http.ServeFile(writer, req, "app/html/index.html") //pass in user here if it is not nil, so that it can say welcome user!

}

func newgame(writer http.ResponseWriter, req *http.Request) {
	g := game.CreateGame(games)
	games = append(games, g)
	urlParams := url.Values{}
	strGameid := strconv.Itoa(g.Gameid)
	urlParams.Add("gameid", strGameid)
	urlParams.Add("Slot", "-1")
	startGameURL := "/play?" + urlParams.Encode()
	variables := map[string]interface{}{"id": g.Gameid, "p1": g.Player1.Id, "p2": g.Player2.Id, "startGameURL": startGameURL}
	outputHTML(writer, "app/html/newgame.html", variables)
}

func play(writer http.ResponseWriter, req *http.Request) {
	urlVars := req.URL.Query()
	varGameid := urlVars["gameid"][0]
	var outputVars map[string]interface{}
	var human bool
	// g := games[gameid] // This needs to be changed to work with database
	var noPossibleMoves bool
	var gameid, _ = strconv.Atoi(varGameid)
	g := games[gameid]

	player := urlVars["Slot"][0]
	//if there is a move
	if player != "JOE" && player != "STEVE" {
		slot, die, dieIndex, capturePiece := game.ParseVariables(urlVars)
		move := game.MoveType{Slot: slot,
			Die:          die,
			DieIndex:     dieIndex,
			CapturePiece: capturePiece,
		}

		endSlot := move.Slot + move.Die
		endSlotState := g.State[endSlot]
		if game.WillCapturePiece(endSlotState, g.CurrTurn.Color) {
			move.CapturePiece = true
			g.Captured[endSlotState] += 1
		}
		fmt.Printf("player %s chose move %v \n", g.CurrTurn.Color, move)

		g.UpdateDice(dieIndex)
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
	g.UpdateTurn()

	//rolls the dice if the dice list is empty
	if len(g.Dice) == 0 {
		g.Dice = game.RollDice(2)
		fmt.Printf("diceroll: %v \n", g.Dice)
	}

	possibleMoves := g.GetPossibleMoves(g.Dice, g.CurrTurn.Color)

	//deletes all dice if no possible moves
	if len(possibleMoves) == 0 {
		g.Dice = nil
		noPossibleMoves = true
	}

	if g.CurrTurn.Id != "JOE" && g.CurrTurn.Id != "STEVE" {
		fmt.Println("human move now")
		human = true
		var urlList []string
		if len(possibleMoves) == 0 {
			urlParams := url.Values{}
			strValues := game.ConvertParams(-1, 0, 0, false)
			urlParams.Add("gameid", varGameid)
			urlParams.Add("Slot", strValues[0])
			var urlString string = "/play?" + urlParams.Encode()
			urlList = append(urlList, urlString)
		} else {
			for index, move := range possibleMoves {
				_ = index
				urlParams := url.Values{}
				strValues := game.ConvertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
				urlParams.Add("gameid", varGameid)
				urlParams = game.AddUrlParams(urlParams, strValues)
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
			move := game.GetAIMove(possibleMoves, g.CurrTurn.Color)
			strValues := game.ConvertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
			urlParams = game.AddUrlParams(urlParams, strValues)
		} else {
			urlParams.Add("Slot", "-1")
		}
		url := "/play?" + urlParams.Encode()
		outputVars = map[string]interface{}{"url": url, "isHuman": human, "state": g.State, "captured": g.Captured, "player": g.CurrTurn.Id}
	}
	games[gameid] = g
	outputHTML(writer, "app/html/playing.html", outputVars)
}

// todo: if someone has won, update the database with wins/losses for each player. Print final board.
func won(writer http.ResponseWriter, req *http.Request) {
	winner := req.URL.Query().Get("winner")
	variables := map[string]interface{}{"winner": winner}
	outputHTML(writer, "app/html/won.html", variables)
}

// just for fun (?), we try to establish a connection to the database here
func dbHandler(writer http.ResponseWriter, req *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Print("successfully connected to database")
}

func initDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Print("successfully connected to database")
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", help) //this makes an endpoint that calls the help function
	http.HandleFunc("/newgame", newgame)
	http.HandleFunc("/play", play)
	//http.HandleFunc("/testplay", testplay)
	http.HandleFunc("/login", login)
	http.HandleFunc("/registered", register)
	http.HandleFunc("/loggedin", loggedin)
	http.HandleFunc("/won", won)
	http.HandleFunc("/db", dbHandler)
	//http.HandleFunc("/scoreboard", scoreboard)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
