package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

// This method send GET and POST request
//and check output

// /user - get all users
// /user/{user_name} - get only this user
// /user - post create new column about user on db
// put - re:check information and writing in file
// delete - delete all/* user(s) info

type Worker struct {
	JobsChan chan string
	SafeZone SafeMapState
}

type SafeMapState struct {
	Mu     sync.Mutex
	Status map[string]*OutputData
}

// Check all users
func (j *Worker) getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var allUser []string
	for _, value := range j.SafeZone.Status {
		allUser = append(allUser, value.Output.Username)
	}
	json.NewEncoder(w).Encode(allUser)
}

// Info about user
func (j *Worker) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//
	representation := make(map[string]string)
	checkAvailability := j.SafeZone.Status
	if len(checkAvailability) == 0 {
		http.Error(w, "Incorrect id: Please search real data:", http.StatusNotFound)
		return
	}

	for key, _ := range j.SafeZone.Status {
		if key == params["id"] {
			representation[params["id"]] = checkAvailability[params["id"]].Output.Username
			err := json.NewEncoder(w).Encode(representation[params["id"]])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			return
		} else {
			continue
		}
	}
	http.Error(w, "Incorrect id: Please search real data:", http.StatusNotFound)
	return
}

func (j *Worker) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//Check ID
	checkAvailability := j.SafeZone.Status
	if len(checkAvailability) == 0 {
		http.Error(w, "Incorrect id: Please search real data:", http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(checkAvailability[params["id"]])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for key, _ := range j.SafeZone.Status {
		if key == params["id"] {
			err := json.NewEncoder(w).Encode(checkAvailability[params["id"]])
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			return
		} else {
			continue
		}
	}
	http.Error(w, "Incorrect id: Please search real data:", http.StatusNotFound)
	return
}

// Create user
func (j *Worker) ParseUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strconv.Itoa(rand.Intn(1000000))
	j.SafeZone.Status[id] = new(OutputData)
	err := json.NewDecoder(r.Body).Decode(&j.SafeZone.Status[id].Output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// empty map
	for key, value := range j.SafeZone.Status {
		if value.Output.Username == j.SafeZone.Status[id].Output.Username && key != id {
			// return if we have this value in map
			j.SafeZone.Status[id].Output.Username = ""
			err = json.NewEncoder(w).Encode(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			j.JobsChan <- j.SafeZone.Status[id].Output.Username
			j.Create(id)
		}
	}

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Update user
/*func (j *Worker) PutUser(w http.ResponseWriter, r *http.Request) {
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

// Delete user/s
func (j *Worker) DeleteUser(w http.ResponseWriter, r *http.Request) {
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


*/

func (j *Worker) Update(id string, ii ImportantInfo) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id] = &OutputData{
		State:  "Done",
		Output: ii,
	}
}

func (j *Worker) Create(id string) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id] = &OutputData{
		State:  "On Work",
		Output: ImportantInfo{Username: j.SafeZone.Status[id].Output.Username},
	}

}

func (j *Worker) GetErrStatus(id string, errCode string) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id] = &OutputData{
		State: errCode,
	}
}

func New(bufferSize int) *Worker {
	// Буферизация на 10 элементов
	jobs := make(chan string, bufferSize)
	status := make(map[string]*OutputData)
	w := &Worker{JobsChan: jobs, SafeZone: SafeMapState{Status: status}}
	return w
}

// Не очень понятно зачем так Но если надо То ок
// Если не сложно поясни пожалуйста
func (j *Worker) Start() {
	ConnectionAPI(j)
}

func ConnectionAPI(w *Worker) {
	log.Info("Server listen ...")
	r := mux.NewRouter()
	r.HandleFunc("/users", w.getAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", w.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}/status", w.GetUserStatus).Methods("GET")
	r.HandleFunc("/users", w.ParseUser).Methods("POST")
	//r.HandleFunc("/users/{id}", w.PutUser).Methods("PUT")
	//r.HandleFunc("/users/{id}", w.DeleteUser).Methods("DELETE")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Error(err)
	}
}
