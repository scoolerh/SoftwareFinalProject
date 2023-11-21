package main

import (
	"backgammon/game"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

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
var g game.Game
var db *sql.DB
var currentUser = "Guest"

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

// Open the home page
func home(writer http.ResponseWriter, req *http.Request) {
	outputHTML(writer, "app/html/index.html", currentUser)
}

func login(writer http.ResponseWriter, req *http.Request) {
	outputHTML(writer, "app/html/login.html", currentUser)
}

func register(writer http.ResponseWriter, req *http.Request) {
	var username string
	var password string

	if req.Method == http.MethodPost {
		username = strings.ToLower(req.FormValue("username"))
		password = strings.ToLower(req.FormValue("password"))
	}

	var query string
	var err error

	query = "SELECT COUNT(*) FROM users WHERE username='" + username + "'"

	var count int
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		panic(err) //might want to change this later
	}

	if count != 0 {
		message := map[string]interface{}{"message": "username taken"}
		outputHTML(writer, "app/html/loginfailed.html", message)
		return
	} else {

		query = "INSERT INTO users VALUES ('" + username + "', '" + password + "')"
		_, err = db.Exec(query)
		if err != nil {
			log.Printf("Error with query %v. Error: %v", query, err)
			panic(err) //might want to change this later
		}

		query = "INSERT INTO userstats (username, gamesPlayed, wins, losses) VALUES ('" + username + "', 0, 0, 0);"
		_, err = db.Exec(query)
		if err != nil {
			log.Printf("Error with query %v. Error: %v", query, err)
			panic(err) //might want to change this later
		}

		currentUser = username
		outputHTML(writer, "app/html/index.html", currentUser)
	}
}

func loggedin(writer http.ResponseWriter, req *http.Request) {
	var username string
	var password string

	if req.Method == http.MethodPost {
		username = strings.ToLower(req.FormValue("username"))
		password = strings.ToLower(req.FormValue("password"))
	}

	validation := validateLogin(username, password)

	if validation == "invalid user" {
		message := map[string]interface{}{"message": "invalid user"}
		outputHTML(writer, "app/html/loginfailed.html", message)
	} else if validation == "wrong password" {
		message := map[string]interface{}{"message": "wrong password"}
		outputHTML(writer, "app/html/loginfailed.html", message)
	} else if validation == "valid" {
		currentUser = username
		outputHTML(writer, "app/html/index.html", currentUser)
	}
}

func validateLogin(username string, password string) string {
	query := "SELECT password FROM users WHERE username='" + username + "'"

	rows, err := db.Query(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	var refPassword string

	if rows.Next() {
		err = rows.Scan(&refPassword)
		if err != nil {
			panic(err)
		}

		if username == "steve" || username == "joe" || username == "guest" || username == "guest2" {
			return "invalid user"
		} else if password != refPassword {
			return "wrong password"
		} else {
			return "valid"
		}
	} else {
		return "invalid user"
	}
}

func selectPlayers(writer http.ResponseWriter, req *http.Request) {
	outputHTML(writer, "app/html/selectPlayers.html", currentUser)
}

func newgame(writer http.ResponseWriter, req *http.Request) {
	var initialState [26]string

	urlVars := req.URL.Query()
	p1 := urlVars["player1"][0]
	p2 := urlVars["player2"][0]

	if p1 == "guest" && p2 == "guest" {
		p1 = "guest"
		p2 = "guest2"
	}

	if p1 == "loggedUser" {
		p1 = currentUser
	}
	if p2 == "loggedUser" {
		p2 = currentUser
	}

	loginmessage := ""
	if p2 == "friend" {
		password := urlVars["password"][0]
		username := urlVars["username"][0]
		if p1 == username {
			loginmessage = "self login"
			p2 = "guest"
		} else {
			validation := validateLogin(username, password)
			if validation == "valid" {
				p2 = username
			} else {
				loginmessage = "invalid login"
				p2 = "guest"
			}
		}
	}

	g, initialState = game.CreateGame(games, p1, p2)

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

	query := "INSERT INTO games (white, black, status, boardstate) VALUES ('" + white + "', '" + black + "', 'new', '" + state + "') RETURNING gameId"

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	var gameid string
	for rows.Next() {
		rows.Scan(&gameid)
	}
	g.Gameid = gameid
	games = append(games, g)

	variables := map[string]interface{}{"login": loginmessage, "currentUser": currentUser, "id": g.Gameid, "p1": g.Player1.Id, "p2": g.Player2.Id}
	outputHTML(writer, "app/html/newgame.html", variables)
}

func rollToStart(writer http.ResponseWriter, req *http.Request) {
	var starter string
	g.Dice = game.RollDice(2)
	for g.Dice[0] == g.Dice[1] {
		g.Dice = game.RollDice(2)
	}

	if g.Dice[0] > g.Dice[1] {
		g.CurrTurn = g.Player1
		starter = g.Player1.Id

	} else {
		g.CurrTurn = g.Player2
		starter = g.Player2.Id
	}

	urlParams := url.Values{}
	urlParams.Add("gameid", g.Gameid)
	urlParams.Add("Slot", "-1")
	startGameURL := "/play?" + urlParams.Encode()
	variables := map[string]interface{}{"currentUser": currentUser, "starter": starter, "die1": g.Dice[0], "die2": g.Dice[1], "id": g.Gameid, "p1": g.Player1.Id, "p2": g.Player2.Id, "startGameURL": startGameURL, "one": g.State[1], "two": g.State[2], "three": g.State[3], "four": g.State[4], "five": g.State[5], "six": g.State[6], "seven": g.State[7], "eight": g.State[8], "nine": g.State[9], "ten": g.State[10], "eleven": g.State[11], "twelve": g.State[12], "thirteen": g.State[13], "fourteen": g.State[14], "fifteen": g.State[15], "sixteen": g.State[16], "seventeen": g.State[17], "eighteen": g.State[18], "nineteen": g.State[19], "twenty": g.State[20], "twentyone": g.State[21], "twentytwo": g.State[22], "twentythree": g.State[23], "twentyfour": g.State[24], "whitehome": g.State[25], "blackhome": g.State[0]}
	outputHTML(writer, "app/html/rollToStart.html", variables)
}

func play(writer http.ResponseWriter, req *http.Request) {
	newRoll := false
	urlVars := req.URL.Query()
	gameid := urlVars["gameid"][0]
	var outputVars map[string]interface{}
	var human bool
	var noPossibleMoves bool

	//if there is a move
	if urlVars["Slot"][0] != "-1" {
		move := parseVariables(urlVars)
		endSlot := move.Slot + move.Die
		endSlotState := g.State[endSlot]
		if game.WillCapturePiece(endSlotState, g.CurrTurn.Color) {
			move.CapturePiece = true
			g.Captured[endSlotState] += 1
		}
		g.UpdateDice(move.DieIndex)
		g.UpdateCaptured(move)
		g.UpdateState(g.CurrTurn.Color, move)
		if g.IsWon() != "" {
			winner := g.CurrTurn.Id
			http.Redirect(writer, req, "/won?winner="+winner, http.StatusSeeOther)
		}

	}

	//rolls the dice if the dice list is empty
	if len(g.Dice) == 0 {
		g.UpdateTurn()
		newRoll = true
		g.Dice = game.RollDice(2)
	}

	possibleMoves := g.GetPossibleMoves(g.Dice, g.CurrTurn.Color)

	//deletes all dice if no possible moves
	if len(possibleMoves) == 0 {
		g.Dice = nil
		noPossibleMoves = true
	}

	if g.CurrTurn.Id != "joe" && g.CurrTurn.Id != "steve" {
		human = true
		urlList, moveList := makeHumanURLs(possibleMoves, gameid)
		if newRoll {
			newRollURL := makeNewRollURL(g.Gameid)
			outputVars = map[string]interface{}{"newRollURL": newRollURL, "game": g, "newRoll": newRoll, "isHuman": human, "noPossibleMoves": noPossibleMoves, "state": g.State, "captured": g.Captured, "player": g.CurrTurn.Id, "one": g.State[1], "two": g.State[2], "three": g.State[3], "four": g.State[4], "five": g.State[5], "six": g.State[6], "seven": g.State[7], "eight": g.State[8], "nine": g.State[9], "ten": g.State[10], "eleven": g.State[11], "twelve": g.State[12], "thirteen": g.State[13], "fourteen": g.State[14], "fifteen": g.State[15], "sixteen": g.State[16], "seventeen": g.State[17], "eighteen": g.State[18], "nineteen": g.State[19], "twenty": g.State[20], "twentyone": g.State[21], "twentytwo": g.State[22], "twentythree": g.State[23], "twentyfour": g.State[24], "whitehome": g.State[25], "blackhome": g.State[0]}
		} else {
			outputVars = map[string]interface{}{"possibleMoves": possibleMoves, "urlList": urlList, "movelist": moveList, "dice": g.Dice, "game": g, "isHuman": human, "noPossibleMoves": noPossibleMoves, "state": g.State, "captured": g.Captured, "player": g.CurrTurn.Id, "one": g.State[1], "two": g.State[2], "three": g.State[3], "four": g.State[4], "five": g.State[5], "six": g.State[6], "seven": g.State[7], "eight": g.State[8], "nine": g.State[9], "ten": g.State[10], "eleven": g.State[11], "twelve": g.State[12], "thirteen": g.State[13], "fourteen": g.State[14], "fifteen": g.State[15], "sixteen": g.State[16], "seventeen": g.State[17], "eighteen": g.State[18], "nineteen": g.State[19], "twenty": g.State[20], "twentyone": g.State[21], "twentytwo": g.State[22], "twentythree": g.State[23], "twentyfour": g.State[24], "whitehome": g.State[25], "blackhome": g.State[0]}
		}
	} else {
		human = false
		url := makeAiURL(possibleMoves, gameid)
		outputVars = map[string]interface{}{"url": url, "isHuman": human, "state": g.State, "captured": g.Captured, "player": g.CurrTurn.Id, "one": g.State[1], "two": g.State[2], "three": g.State[3], "four": g.State[4], "five": g.State[5], "six": g.State[6], "seven": g.State[7], "eight": g.State[8], "nine": g.State[9], "ten": g.State[10], "eleven": g.State[11], "twelve": g.State[12], "thirteen": g.State[13], "fourteen": g.State[14], "fifteen": g.State[15], "sixteen": g.State[16], "seventeen": g.State[17], "eighteen": g.State[18], "nineteen": g.State[19], "twenty": g.State[20], "twentyone": g.State[21], "twentytwo": g.State[22], "twentythree": g.State[23], "twentyfour": g.State[24], "whitehome": g.State[25], "blackhome": g.State[0]}
	}
	outputHTML(writer, "app/html/playing.html", outputVars)
}

func won(writer http.ResponseWriter, req *http.Request) {
	winner := req.URL.Query().Get("winner")
	variables := map[string]interface{}{"winner": winner}
	var err error
	var query string
	var loser string

	query = "UPDATE userstats SET gamesPlayed = gamesPlayed + 1 WHERE username = '" + g.Player1.Id + "';"
	//exec here
	_, err = db.Exec(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	query = "UPDATE userstats SET gamesPlayed = gamesPlayed + 1 WHERE username = '" + g.Player2.Id + "';"
	_, err = db.Exec(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	query = "UPDATE userstats SET wins = wins + 1 WHERE username = '" + winner + "';"
	_, err = db.Exec(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	if g.Player1.Id == winner {
		loser = g.Player2.Id
	} else {
		loser = g.Player1.Id
	}

	query = "UPDATE userstats SET losses = losses + 1 WHERE username = '" + loser + "';"
	_, err = db.Exec(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	query = "UPDATE games SET status = 'finished', winner = '" + winner + "' WHERE gameId = " + g.Gameid + ";"
	_, err = db.Exec(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	outputHTML(writer, "app/html/won.html", variables)
}

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
	log.Print("follow this link to play backgammon: http://localhost:9000/")

}

func scoreboard(writer http.ResponseWriter, req *http.Request) {
	pythonScript := "app/db_api.py"
	// Command to run Python script
	cmd := exec.Command("python3", pythonScript)
	outputVars, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing Python script:", err)
		fmt.Println("Python script output:", string(outputVars))
		return
	}
	jsonStr := string(outputVars)

	outputHTML(writer, "app/html/scoreboard.html", jsonStr)

}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", home)
	http.HandleFunc("/selectPlayers", selectPlayers)
	http.HandleFunc("/newgame", newgame)
	http.HandleFunc("/rollToStart", rollToStart)
	http.HandleFunc("/play", play)
	http.HandleFunc("/login", login)
	http.HandleFunc("/registered", register)
	http.HandleFunc("/loggedin", loggedin)
	http.HandleFunc("/scoreboard", scoreboard)
	http.HandleFunc("/won", won)
	http.HandleFunc("/db", dbHandler)
	fs := http.FileServer(http.Dir("app/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
