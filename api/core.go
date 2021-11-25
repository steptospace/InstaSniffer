package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

// This method send GET and POST request
//and check output

// /user - get all users
// /user/{user_name} - get only this user
// /user - post create new column about user on db
// put - re:check information and writing in file
// delete - delete all/* user(s) info

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
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			//404 прокинуть
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

func (j Worker) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range users {
		if item.Id == params["id"] {
			err := json.NewEncoder(w).Encode(j.SafeZone.Status[item.Name])
			if err != nil {
				// 404?
				fmt.Println("We are at here")
				//log.Error(err)
			}
			return
		}
	}
}

// Create user
func (j Worker) PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Check this one
		j.GetErrStatus(user, "Error: 400")
		return
	}
	users = append(users, user)
	j.JobsChan <- user.Name
	j.Create(user)
}

// Update user
func (j Worker) PutUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.Id == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var book User
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.Id = params["id"]
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
		if item.Id == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}

type OutputData struct {
	State  string
	Output ImportantInfo
	Error  *ErrInfo
}

type Worker struct {
	JobsChan chan string
	SafeZone SafeMapState
}

type SafeMapState struct {
	Mu     sync.Mutex
	Status map[string]*OutputData
}

func (j Worker) Update(username string, state map[string]*OutputData, ii ImportantInfo) {
	j.SafeZone.Mu.Lock()
	tt := state[username]
	tt.State = "Done"
	tt.Output = ii
	j.SafeZone.Mu.Unlock()
}

func (j Worker) Create(user User) {
	j.SafeZone.Mu.Lock()
	j.SafeZone.Status[user.Name] = &OutputData{
		State: "On Work",
	}
	j.SafeZone.Mu.Unlock()
}

func (j Worker) GetErrStatus(user User, errCode string) {
	j.SafeZone.Mu.Lock()
	//Изменить код респонза
	j.SafeZone.Status[user.Name] = &OutputData{
		State: errCode,
	}
	j.SafeZone.Mu.Unlock()
}

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
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Error(err)
	}
}
