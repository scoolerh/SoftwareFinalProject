package main

import (
	"net/url"
	"strconv"
)

// adds the url values to an array
func AddUrlParams(urlParams url.Values, valuesToAdd [4]string) url.Values {
	urlParams.Add("Slot", valuesToAdd[0])
	urlParams.Add("Die", valuesToAdd[1])
	urlParams.Add("DieIndex", valuesToAdd[2])
	urlParams.Add("CapturePiece", valuesToAdd[3])
	return urlParams
}

// converts the parameters that will be inserted into the URL from their respective types into strings
func ConvertParams(slot int, die int, index int, capture bool) [4]string {
	strSlot := strconv.Itoa(slot)
	strDie := strconv.Itoa(die)
	strDieIndex := strconv.Itoa(index)
	strCapturePiece := strconv.FormatBool(capture)
	returns := [4]string{strSlot, strDie, strDieIndex, strCapturePiece}
	return returns
}

// takes the variables from the url and converts to their necessary types
func ParseVariables(urlVariables url.Values) (int, int, int, bool) {
	varSlot := urlVariables["Slot"][0]
	varDie := urlVariables["Die"][0]
	varDieIndex := urlVariables["DieIndex"][0]
	varCapturePiece := urlVariables["CapturePiece"][0]

	slot, _ := strconv.Atoi(varSlot)
	die, _ := strconv.Atoi(varDie)
	dieIndex, _ := strconv.Atoi(varDieIndex)
	capturePiece, _ := strconv.ParseBool(varCapturePiece)
	return slot, die, dieIndex, capturePiece
}
