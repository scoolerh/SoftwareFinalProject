package main

import (
	"backgammon/game"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

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
var initialState = [26]string{"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}
var testState = [26]string{"", "ww", "bb", "w", "b", "ww", "bb", "w", "b", "ww", "bb", "w", "", "", "", "", "b", "ww", "bb", "w", "b", "ww", "bb", "w", "b", ""}
var p1 game.Player
var p2 game.Player
var whoseTurn string = "first"
var gameid int
var winner string
var db *sql.DB

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
	http.ServeFile(writer, req, "./html/index.html")
}

// todo: Create a database for users, allow a user to log in (or sign up if they do not have a username)
func login(writer http.ResponseWriter, req *http.Request) {
	fmt.Printf("Connecting to login endpoint")
	http.ServeFile(writer, req, "/html/login.html")
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

}

// Starts a new game for the user and displays the initial board
func newgame(writer http.ResponseWriter, req *http.Request) {

	p1, p2 = game.Player{Id: "STEVE", Color: "w"}, game.Player{Id: "JOE", Color: "b"} //will need to be an input in the future
	gameid = len(games)
	capturedMap := initializeCapturedMap()
	g := game.Game{Gameid: gameid, Player1: p1, Player2: p2, State: initialState, Captured: capturedMap}
	games = append(games, g)

	var white string
	var black string
	var state string
	if g.Player1.Color == "w" {
		white = g.Player1.Id
		black = g.Player2.Id
	} else {
		white = g.Player2.Id
		black = g.Player1.Id
	}

	for index, slot := range initialState {
		_ = index
		state += slot + "o"
	}

	query := "INSERT INTO games (white, black, status, boardstate) VALUES ('" + white + "', '" + black + "', 'new', '" + state + "')"
	//might do something like this instead to prevent injection
	//func buildSql(email string) string {
	//return fmt.Sprintf("SELECT * FROM users WHERE email='%s';", email)

	var err error
	_, err = db.Exec(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	//TODO: set up db connection here instead. Think about players and gameID
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

		// returning and printing boardState for testing purposes
		log.Printf("\n move nr %v: \n", i)
		fmt.Fprintf(writer, "move nr %v: \n", i)
		g.Move(g.Player1)
		//fmt.Fprintf(writer, "player 1 made a move: %v", g.currMove)
		fmt.Fprintf(writer, "%v", g.State)
		fmt.Fprintf(writer, "captured pieces: %v \n", g.Captured)
		g.Move(g.Player2)
		//fmt.Fprintf(writer, "player 2 made a move: %v", g.currMove)
		fmt.Fprintf(writer, "%v", g.State)
		fmt.Fprintf(writer, "captured pieces: %v \n", g.Captured)
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
		//TODO: update both game and user stats here
		return
	}
	if whoseTurn == "first" {
		variables := map[string]interface{}{"id": gameid, "player": p1.Id, "state": g.State, "captured": g.Captured}
		outputHTML(writer, "./html/playing.html", variables)
		whoseTurn = "w"
		games[gameid] = g
	} else if whoseTurn == "w" {
		variables := map[string]interface{}{"id": gameid, "player": p2.Id, "state": g.State, "captured": g.Captured}
		g.Move(g.Player1)
		outputHTML(writer, "./html/playing.html", variables)
		whoseTurn = "b"
		games[gameid] = g
	} else {
		variables := map[string]interface{}{"id": gameid, "player": p1.Id, "state": g.State, "captured": g.Captured}
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
	http.HandleFunc("/testplay", testplay)
	http.HandleFunc("/login", login)
	http.HandleFunc("/registered", register)
	http.HandleFunc("/won", won)
	http.HandleFunc("/db", dbHandler)
	//http.HandleFunc("/scoreboard", scoreboard)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
