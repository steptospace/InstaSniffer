package API

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
)

// This method send GET and POST request
//and check output

//Осталось определиться что делает get, а что будет делать пост
// По сути нам просто нужно отправлять данные на принимающий узел

// Изначальная идея:
// /user - get all users
// /user/{user_name} - get only this user
// /user - post create new column about user on db
// put - re:check information and writing in file
// delete - delete all/* user(s) info

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users []User

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range users {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})

}

// Create user (Test mode)
func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = strconv.Itoa(rand.Intn(1000000))
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

//Update user
func PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var book User
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			users = append(users, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}

func ConnectionAPI() {
	fmt.Println("Hello are u here Docker?")
	r := mux.NewRouter()
	users = append(users, User{ID: "1", Name: "Jonny"})
	users = append(users, User{ID: "2", Name: "Lilly"})
	r.HandleFunc("/users", getAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users", PostUser).Methods("POST")
	r.HandleFunc("/users/{id}", PutUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	log.Error(http.ListenAndServe(":8000", r))
}
