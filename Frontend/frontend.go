package main

import "fmt"


func main() {
	// initialState := map[string]string{"0": "", "1": "ww", "2": "", "3": "", "4": "", "5": "", "6": "bbbbbb", "7": "", "8": "bbb", "9": "", "10": "", "11": "", "12": "wwwwww",
	// "13": "bbbbb", "14": "", "15": "", "16": "", "17": "www", "18": "", "19": "wwwww", "20": "", "21": "", "22": "", "23": "", "24": "bb", "25": ""}

	initialState := [26]string {"", "ww", "", "", "", "", "bbbbb", "", "bbb", "", "", "", "wwwww", "bbbbb", "", "", "", "www", "", "wwwww", "", "", "", "", "bb", ""}
	fmt.Println(initialState)

	for spot, pieces := range initialState {
		fmt.Println(spot, pieces)
	}

	// for position, pieces := range initialState {
	// 	fmt.Println(position)
	// 	fmt.Println(pieces)
	// }
}

