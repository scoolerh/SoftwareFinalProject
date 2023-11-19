package main

import (
	"backgammon/game"
	"net/url"
	"strconv"
)

// adds the url values to an array
func addUrlParams(urlParams url.Values, valuesToAdd [4]string) url.Values {
	urlParams.Add("Slot", valuesToAdd[0])
	urlParams.Add("Die", valuesToAdd[1])
	urlParams.Add("DieIndex", valuesToAdd[2])
	urlParams.Add("CapturePiece", valuesToAdd[3])
	return urlParams
}

// creates a url for the AI and gets the AI's move when necessary
func makeAiURL(possibleMoves []game.MoveType, gameid string) string {
	urlParams := url.Values{}
	urlParams.Add("gameid", gameid)
	if len(possibleMoves) != 0 {
		move := game.GetAIMove(possibleMoves, g.CurrTurn, g)
		strValues := convertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
		urlParams = addUrlParams(urlParams, strValues)
	} else {
		urlParams.Add("Slot", "-1")
	}
	url := "/play?" + urlParams.Encode()
	return url
}

// creates a url for when the human has no move to make
func makeHumanNoMoveURL(gameid string) ([]string, [][3]int) {
	var urlList []string
	var moveList [][3]int
	urlParams := url.Values{}
	urlParams.Add("gameid", gameid)
	urlParams.Add("Slot", "-1")
	var urlString string = "/play?" + urlParams.Encode()
	urlList = append(urlList, urlString)
	move := [3]int{0, 0, 0}
	moveList = append(moveList, move)
	return urlList, moveList
}

// creates the urls when the human has move(s) to make
func makeHumanMovesUrls(possibleMoves []game.MoveType, gameid string) ([]string, [][3]int) {
	var urlList []string
	var moveList [][3]int
	for _, move := range possibleMoves {
		urlParams := url.Values{}
		strValues := convertParams(move.Slot, move.Die, move.DieIndex, move.CapturePiece)
		urlParams.Add("gameid", gameid)
		urlParams = addUrlParams(urlParams, strValues)
		var urlString string = "/play?" + urlParams.Encode()
		urlList = append(urlList, urlString)
		move := [3]int{move.Slot, move.Slot + move.Die, move.Die}
		moveList = append(moveList, move)
	}
	return urlList, moveList
}

func makeHumanURLs(possibleMoves []game.MoveType, gameid string) ([]string, [][3]int) {
	var urlList []string
	var moveList [][3]int
	if len(possibleMoves) == 0 {
		urlList, moveList = makeHumanNoMoveURL(gameid)
	} else {
		urlList, moveList = makeHumanMovesUrls(possibleMoves, gameid)
	}
	return urlList, moveList
}

// creates the url for when there is a new roll
func makeNewRollURL(gameid string) string {
	urlParams := url.Values{}
	urlParams.Add("gameid", gameid)
	urlParams.Add("Slot", "-1")
	newRollURL := "/play?" + urlParams.Encode()
	return newRollURL
}

// converts the parameters that will be inserted into the URL from their respective types into strings
func convertParams(slot int, die int, index int, capture bool) [4]string {
	strSlot := strconv.Itoa(slot)
	strDie := strconv.Itoa(die)
	strDieIndex := strconv.Itoa(index)
	strCapturePiece := strconv.FormatBool(capture)
	returns := [4]string{strSlot, strDie, strDieIndex, strCapturePiece}
	return returns
}

// takes the variables from the url and converts to their necessary types
func parseVariables(urlVariables url.Values) game.MoveType {
	varSlot := urlVariables["Slot"][0]
	varDie := urlVariables["Die"][0]
	varDieIndex := urlVariables["DieIndex"][0]
	varCapturePiece := urlVariables["CapturePiece"][0]

	slot, _ := strconv.Atoi(varSlot)
	die, _ := strconv.Atoi(varDie)
	dieIndex, _ := strconv.Atoi(varDieIndex)
	capturePiece, _ := strconv.ParseBool(varCapturePiece)

	move := game.MoveType{Slot: slot,
		Die:          die,
		DieIndex:     dieIndex,
		CapturePiece: capturePiece,
	}
	return move
}
