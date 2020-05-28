package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// what should happend?
// i should seperate it into two basic layers:
// 1) The first one is getting data and constructing two arrays of 'points'
// (that is points in meters, not coordinates)
// 2) The second one is to take these coordinates and produce a "result"
//

func main() {
	/*
		fmt.Println("Hello")
		shapeA := []point{point{0, 0}, point{1, 1}, point{2, 4}}
		shapeB := []point{point{0, 0}, point{1, 1}, point{3, 4}}
		test := difference(shapeA, shapeB)
		fmt.Println(test)
		sortInput("oh no")
		fmt.Println("After")
	*/

	http.HandleFunc("/", handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3100"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type person struct {
	Name string
	Age  int
}

type someThing struct {
	One string
	Two string
}

type testJSON struct {
	Test  string
	Thing someThing
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("")
	fmt.Println("<<<SESSION START>>>")

	var input request
	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err2 := json.NewDecoder(r.Body).Decode(&input)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
		println("Could not parse json")
		return
	}

	result, err := response(input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Could not validate")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

	fmt.Println("<<<SESSION END>>>")
	fmt.Println("")

}

func response(req request) (swirlDecision, error) {
	if !valdiate(req) {
		return swirlDecision{}, errors.New("The request could not be validated")
	}
	return swirlRunsDeep(req.RunInput.RunA, req.RunInput.RunB), nil
}

/*
// "coors": [{"lat": 30.0, "long": 1.0}, {"lat": 31.3, "long": 2.0}, {"lat": 32.3, "long": 3.0}]
	// "coors": [[30.0, 1.0], [31.3, 2.0], [32.3, 3.0]]
}	"coors": [1, 2, 3]
*/
