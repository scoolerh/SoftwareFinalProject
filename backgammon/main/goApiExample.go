package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func plus(x int, y int) int {
	return x + y
}

func minus(x int, y int) int {
	return x - y
}

func math(writer http.ResponseWriter, req *http.Request) {

	xString := req.URL.Query().Get("x")
	yString := req.URL.Query().Get("y")

	x, errX := strconv.Atoi(xString)
	y, errY := strconv.Atoi(yString)

	// if errX != nil {
	// 	fmt.Fprint(writer, "Error occurred")
	// 	log.Fatal(errX)
	// }

	// if errY != nil {
	// 	fmt.Fprint(writer, "Error occurred")
	// 	log.Fatal(errY)
	// }

	if errX != nil || errY != nil {
		errorMessage := "Invalid input for x and/or y"
		http.Error(writer, fmt.Sprintf("%s - %s", errorMessage, http.StatusText(http.StatusBadRequest)), http.StatusBadRequest)
		return
	}

	var mathFunction func(int, int) int

	if x > 5 {
		mathFunction = plus
	} else {
		mathFunction = minus
	}

	z := mathFunction(x, y)

	fmt.Fprint(writer, "This is a math page \n")

	fmt.Fprintf(writer, "The answer is %d", z)

}

func exampleMain() {
	http.HandleFunc("/math", math)
	http.ListenAndServe(":5555", nil) //listens for HTTP on port 5555, with standard mapping
}
