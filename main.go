package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	// add `json` to set prop name when convert to json
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	LastName string  `json:"lastName"`
	Age      int     `json:"age"`
	Height   float32 `json:"height"`
}

var Users []User = []User{
	User{
		Id:       1,
		Name:     "John",
		LastName: "Locke",
		Age:      45,
		Height:   1.75,
	},
	User{
		Id:       2,
		Name:     "Jack",
		LastName: "Shephard",
		Age:      40,
		Height:   1.72,
	},
	User{
		Id:       3,
		Name:     "James",
		LastName: "Ford",
		Age:      39,
		Height:   1.77,
	},
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Users)
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//change status to 201: created
	w.WriteHeader(http.StatusCreated)

	body, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Fprintf(w, "Error on postUsers method")
		return
	}
	var newUser User
	json.Unmarshal(body, &newUser)
	newUser.Id = len(Users) + 1
	Users = append(Users, newUser)
	json.NewEncoder(w).Encode(newUser)
}

func handlerRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", getUsers).Methods("GET")
	myRouter.HandleFunc("/users", postUsers).Methods("POST")
	log.Fatal(http.ListenAndServe(":3003", myRouter))
}

func main() {
	fmt.Println("Server running on the port 3003")
	handlerRequests()
}
