package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var users = map[string]string{
	"Alex": "abcdefg",
}

func valdiate(req request) bool {
	name := req.User.Name
	pswd := req.User.Password
	if users[name] != pswd {
		return false
	}
	if len(req.RunInput.RunA.Coors) < 3 || len(req.RunInput.RunB.Coors) < 3 {
		return false
	}
	return true
}

// for debugging etc

func printJSON(points []point) {
	b, err := json.Marshal(&points)
	if err != nil {
		fmt.Println("Could not make points to json")
		return
	}
	// fmt.Println(string(b))
	f, err := os.Create("tmp/curve.json")
	defer f.Close()
	f.WriteString(string(b))
}
