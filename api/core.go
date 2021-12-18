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
	Mu     sync.RWMutex
	Status map[string]*OutputData
}

// Check all users
func (j *Worker) getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
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
	data, notFound := j.getUserInfo(params["id"])
	if notFound {
		//Set params in new structure User can see json request and info about err in sys
		j.SetErrStatus(params["id"], http.StatusNotFound, "Cant find this id Please check again")
		http.Error(w, "Cant find this id Please check id again", http.StatusNotFound)
		return
	} else {
		err := json.NewEncoder(w).Encode(data.Output.Username)
		if err != nil {
			j.SetErrStatus(params["id"], http.StatusBadRequest, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}
}

func (j *Worker) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// 1. getUserInfo() (возвращает OutoutData, bool) (используется mutex Rlock)
	// 2. if notFound createError(id) (используется mutex lock)
	// 3. else json.NewEncoder(w).Encode(OutoutData)
	data, notFound := j.getUserInfo(params["id"])
	if notFound {
		//Set params in new structure User can see json request and info about err in sys
		j.SetErrStatus(params["id"], http.StatusNotFound, "Cant find this id Please check again")
		http.Error(w, "Cant find this id Please check id again", http.StatusNotFound)
		return
	} else {
		err := json.NewEncoder(w).Encode(data)
		// Нужна ли эта проверка Если честно я запутался
		if err != nil {
			j.SetErrStatus(params["id"], http.StatusBadRequest, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}
}

// Create user
func (j *Worker) ParseUser(w http.ResponseWriter, r *http.Request) {
	// Не уверен что тут подойдет создание
	var task User
	var user Inside
	w.Header().Set("Content-Type", "application/json")
	user.Id = strconv.Itoa(rand.Intn(1000000))
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		j.SetErrStatus(user.Id, http.StatusBadRequest, "We cant decode this data Check our params")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Username == "" {
		j.SetErrStatus(user.Id, http.StatusBadRequest, "Incorrect Username please check this value and try again")
		http.Error(w, "Empty username params", http.StatusBadRequest)
		return
	}

	key, isDup := j.isDuplicate(task.Username)
	if isDup {
		err = json.NewEncoder(w).Encode(key)
		if err != nil {
			j.SetErrStatus(user.Id, http.StatusInternalServerError, "Server error try again")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		j.Create(user.Id, task.Username)
		j.JobsChan <- user.Id
	}
	// if task.Username empty return err 400
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		j.SetErrStatus(user.Id, http.StatusInternalServerError, "Server error try again")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//Check duplicate in map
func (j *Worker) isDuplicate(username string) (uid string, isDup bool) {
	// empty map
	j.SafeZone.Mu.RLock()
	defer j.SafeZone.Mu.RUnlock()
	for key, value := range j.SafeZone.Status {
		if username == value.Output.Username {
			isDup = true
			return key, isDup
		}
	}
	return uid, isDup
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

func (j *Worker) Create(id string, username string) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id] = &OutputData{
		State:  "On Work",
		Output: ImportantInfo{Username: username},
	}

}

func (j *Worker) SetErrStatus(id string, errCode int, description string) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id] = &OutputData{
		State: "Error",
		Error: &ErrInfo{Err: errCode, Description: description},
	}
}

func (j *Worker) getUserInfo(id string) (data OutputData, notFound bool) {
	j.SafeZone.Mu.RLock()
	defer j.SafeZone.Mu.RUnlock()
	if _, ok := j.SafeZone.Status[id]; ok {
		// if we find user on system
		data = *j.SafeZone.Status[id]
		return data, notFound
	} else {
		notFound = true
	}
	return data, notFound
}

// naming вери хард фор ми хелп ми плиз
func (j *Worker) GetUsernameById(id string) (username string, isUsed bool) {
	j.SafeZone.Mu.RLock()
	defer j.SafeZone.Mu.RUnlock()
	if _, ok := j.SafeZone.Status[id]; ok {
		isUsed = true
		username = j.SafeZone.Status[id].Output.Username
		return username, isUsed
	} else {
		isUsed = false
	}
	return username, isUsed
}

func New(bufferSize int) *Worker {
	jobs := make(chan string, bufferSize)
	status := make(map[string]*OutputData)
	w := &Worker{JobsChan: jobs, SafeZone: SafeMapState{Status: status}}
	return w
}

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
