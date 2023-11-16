package main

import (
	"backgammon/game"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"log"
	"os/exec"

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

// Open the home page 
func home(writer http.ResponseWriter, req *http.Request) {
	variables := map[string]interface{}{"username": ""}
	outputHTML(writer, "app/html/index.html", variables)
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
	variables := map[string]interface{}{"username": username}
	outputHTML(writer, "app/html/index.html", variables)
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
	variables := map[string]interface{}{"username": currentUser}
	outputHTML(writer, "app/html/index.html", variables)
}

func selectOpponent(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "app/html/selectOpponent.html")
}

func selectAI(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "app/html/selectAI.html")
}

func selectHuman(writer http.ResponseWriter, req *http.Request) {
	http.ServeFile(writer, req, "app/html/selectHuman.html")
}

func newgame(writer http.ResponseWriter, req *http.Request) {
	var initialState [26]string
	g, initialState = game.CreateGame(games)

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
	//might do something like this instead to prevent injection
	//func buildSql(email string) string {
	//return fmt.Sprintf("SELECT * FROM users WHERE email='%s';", email)

	rows, err := db.Query(query)
	if err != nil {
		panic(err) //might want to change this later
	}

	var gameid string
	for rows.Next() {
		rows.Scan(&gameid)
	}
	//is this actually inserting into the struct?
	g.Gameid = gameid
	games = append(games, g)

	urlParams := url.Values{}
	urlParams.Add("gameid", g.Gameid)
	log.Printf("Gameid in urlParams: %v", urlParams["gameid"][0])
	urlParams.Add("Slot", "-1")
	startGameURL := "/play?" + urlParams.Encode()
	log.Printf("startGameURL: %v", startGameURL)
	variables := map[string]interface{}{"id": g.Gameid, "p1": g.Player1.Id, "p2": g.Player2.Id, "startGameURL": startGameURL}
	outputHTML(writer, "app/html/newgame.html", variables)
}

func play(writer http.ResponseWriter, req *http.Request) {
	urlVars := req.URL.Query()
	log.Printf("url: %v", req.URL)
	log.Printf("urlVars: %v", urlVars)
	gameid := urlVars["gameid"][0]
	var outputVars map[string]interface{}
	var human bool
	// g := games[gameid] // This needs to be changed to work with database
	var noPossibleMoves bool
	// var intGameid, _ = strconv.Atoi(gameid)
	// g := games[intGameid]

	//if there is a move
	if urlVars["Slot"][0] != "-1" {
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
			urlParams.Add("gameid", gameid)
			urlParams.Add("Slot", strValues[0])
			var urlString string = "/play?" + urlParams.Encode()
			urlList = append(urlList, urlString)
		} else {
			for index, move := range possibleMoves {
				_ = index
				urlParams := url.Values{}
				strValues := game.ConvertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
				urlParams.Add("gameid", gameid)
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
		urlParams.Add("gameid", gameid)
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
	// games[intGameid] = g
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

func scoreboard(writer http.ResponseWriter, req *http.Request) {
	err := exec.Command("python", "db_api.py").Start()
	if err != nil {
		log.Fatalf("Error getting current directory: %s", err)
	}
	pythonScript := filepath.Join(dir, "db_api.py")

	//make sure python file exists
	if _, err := os.Stat(pythonScript); os.IsNotExist(err) {
		log.Fatalf("Python script '%s' does not exist", pythonScript)
	}
	// Command to run Python script
	cmd := exec.Command("python", pythonScript)
	
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

	http.HandleFunc("/", home) 
	http.HandleFunc("/selectOpponent", selectOpponent) 
	http.HandleFunc("/selectAI", selectAI) 
	http.HandleFunc("/selectHuman", selectHuman) 
	http.HandleFunc("/newgame", newgame)
	http.HandleFunc("/play", play)
	http.HandleFunc("/login", login)
	http.HandleFunc("/registered", register)
	http.HandleFunc("/loggedin", loggedin)
	http.HandleFunc("/won", won)
	http.HandleFunc("/db", dbHandler)
	http.HandleFunc("/scoreboard", scoreboard)
	fs := http.FileServer(http.Dir("app/static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 9000, with standard mapping
}
