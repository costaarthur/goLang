package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	// newEncoder returns a new encoder that writes to w
	// transform Users into json
	// encode writes on the screen
	json.NewEncoder(w).Encode(Users)
}

func filterUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	splitedPath := strings.Split(r.URL.Path, "/")

	if len(splitedPath) > 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//convert splitedPath (string) to int
	id, err := strconv.Atoi(splitedPath[2])
	if err != nil {
		fmt.Fprintf(w, "Error on filtertUsers method")
		return
	}

	for _, user := range Users {
		if user.Id == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Error on postUsers method")
		return
	}
	var newUser User
	//parses the json.encoded data and store the result
	json.Unmarshal(body, &newUser)
	newUser.Id = len(Users) + 1
	Users = append(Users, newUser)
	json.NewEncoder(w).Encode(newUser)
	//change status to 201: created
	w.WriteHeader(http.StatusCreated)
}

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	var userToDelete User
	json.Unmarshal(body, &userToDelete)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for index, user := range Users {
		if userToDelete.Id == user.Id {
			Users = append(Users[0:index], Users[index+1:len(Users)]...)
			w.WriteHeader(http.StatusNoContent)
			break
		}
	}
	fmt.Println("entrei")
	w.WriteHeader(http.StatusNotFound)
}

func handlerRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", getUsers).Methods("GET")
	myRouter.HandleFunc("/users/", filterUsers).Methods("GET")
	myRouter.HandleFunc("/users", postUsers).Methods("POST")
	myRouter.HandleFunc("/users", deleteUsers).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3003", myRouter))
}

func main() {
	fmt.Println("Server running on the port 3003")
	handlerRequests()
}
