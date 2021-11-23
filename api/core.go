package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// This method send GET and POST request
//and check output

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

// Check all users
func (j Worker) getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Information about user
// Output info on web-wall
func (j Worker) GetUser(w http.ResponseWriter, r *http.Request) {
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

func (j Worker) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range users {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(j.Status[item.Name])
			return
		}
	}
}

// Create user
func (j Worker) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
	j.JobsChan <- user.Name
	//
	j.Status[user.Name] = &OutputData{
		State: "On Work",
	}
}

// Update user
func (j Worker) PutUser(w http.ResponseWriter, r *http.Request) {
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

func (j Worker) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

type ImportantInfo struct {
	Username   string
	Name       string
	Bio        string
	Created_at time.Time
	Images     []Media
	Videos     []Media
	Avatar     string // link on photo
}

type Media struct {
	Url         string
	Description string
}

type OutputData struct {
	State  string
	Output ImportantInfo
}

type Worker struct {
	JobsChan chan string
	Status   map[string]*OutputData
}

// Написать функцию записи в файл Функция маршал

func ConnectionAPI(w Worker) {
	fmt.Println("Server listen ...")
	r := mux.NewRouter()
	// examples case

	r.HandleFunc("/users", w.getAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", w.GetUser).Methods("GET")
	//Добавить роут /users/{id}/result. рез статус *(в работе или завершен)
	//В ответе json с полями status и result (ImportantInfo)
	r.HandleFunc("/users/{id}/status", w.GetUserStatus).Methods("GET")
	r.HandleFunc("/users", w.PostUser).Methods("POST")
	r.HandleFunc("/users/{id}", w.PutUser).Methods("PUT")
	r.HandleFunc("/users/{id}", w.DeleteUser).Methods("DELETE")
	log.Error(http.ListenAndServe(":8000", r))
}
