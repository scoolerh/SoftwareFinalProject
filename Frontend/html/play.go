//make buttons for possible moves
//when you click the button, tells the backend the chosen move

import "text/template"

data := pageData{
	gameState: ,
	possibleMoves: GetPossibleMoves,
}
var possibleMoves []MoveType = GetPossibleMoves

<ul>
	{{range .possibleMoves}}
		<button type="button">{{ .Slot .Die}}</button>
	{{end}}
</ul>
