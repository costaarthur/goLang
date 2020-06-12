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
	fmt.Fprintf(w, "Welcome home, sir.")
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

	vars := mux.Vars(r)
	// convert id to number
	id, _ := strconv.Atoi(vars["userId"])

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

func updateUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	splitedPath := strings.Split(r.URL.Path, "/")
	//convert string to int
	id, convertToIntErr := strconv.Atoi(splitedPath[2])

	if convertToIntErr != nil {
		fmt.Fprintf(w, "Error on convert to int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, bodyErr := ioutil.ReadAll(r.Body)

	if bodyErr != nil {
		fmt.Fprintf(w, "Error on read body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var updatedUser User
	jsonErr := json.Unmarshal(body, &updatedUser)

	if jsonErr != nil {
		fmt.Fprintf(w, "Error on convert to json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for index, user := range Users {
		if id == user.Id {
			Users[index] = updatedUser
			json.NewEncoder(w).Encode(updatedUser)
		}
	}
	w.WriteHeader(http.StatusNotFound)
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
	w.WriteHeader(http.StatusNotFound)
}

func handlerRequests(myRouter *mux.Router) {
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/users", getUsers).Methods("GET")
	myRouter.HandleFunc("/users/{userId}", filterUsers).Methods("GET")
	myRouter.HandleFunc("/users", postUsers).Methods("POST")
	myRouter.HandleFunc("/users/", updateUsers).Methods("PUT")
	myRouter.HandleFunc("/users", deleteUsers).Methods("DELETE")
}

func handlerServer() {
	// strictSlash(true) means that the route /users/ redirect to /users, preventing errors
	myRouter := mux.NewRouter().StrictSlash(true)
	handlerRequests(myRouter)
	fmt.Println("Server running on the port 3003")
	log.Fatal(http.ListenAndServe(":3003", myRouter))
}

func main() {
	handlerServer()
}
