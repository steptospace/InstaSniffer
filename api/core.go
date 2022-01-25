package api

import (
	"InstaSniffer/db"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// This method send GET and POST request
//and check output

// /user - get all users
// /user/{user_name} - get only this user
// /user - post create new column about user on db
// put - re:check information and writing in file
// delete - delete all/* user(s) info

const (
	StatusError  = "Error"
	StatusDone   = "Done"
	statusInWork = "In Work"
)

type Worker struct {
	JobsChan chan string
	SafeZone SafeMapState
	DB       *gorm.DB
}

type SafeMapState struct {
	Mu     sync.RWMutex
	Status map[string]*OutputData
}

func (j *Worker) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	var allUser []string
	for _, value := range j.SafeZone.Status {
		allUser = append(allUser, value.Output.Username)
	}
	json.NewEncoder(w).Encode(allUser)
}

func (j *Worker) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//
	data, notFound := j.getUserInfo(params["id"])
	if notFound {
		//Set params in new structure
		//User can see json request and info about err in sys
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrInfo{
			Err:         http.StatusNotFound,
			Description: "Cant find this id Please check again",
		})
		return
	} else {
		if data.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		dbData, err := db.GetInfoById(j.DB, params["id"])
		if err != nil {
			log.Error(err)
		}
		if data.Output.Username == dbData.Name {
			json.NewEncoder(w).Encode(data)
			return
		}
		return
	}
}

func (j *Worker) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	data, notFound := j.getUserInfo(params["id"])
	if notFound {
		//Set params in new structure User can see json request and info about err in sys
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrInfo{
			Err:         http.StatusNotFound,
			Description: "Cant find this id Please check again",
		})
		return
	} else {
		if data.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(data)
		return
	}
}

// Create user
func (j *Worker) ParseUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user Inside
	user.Id = strconv.Itoa(rand.Intn(1000000))

	var task User
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrInfo{
			Err:         http.StatusBadRequest,
			Description: "We cant decode this data Check our params",
		})
		return
	}

	if task.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrInfo{
			Err:         http.StatusBadRequest,
			Description: "Incorrect Username please check this value and try again",
		})
		return
	}

	key, isDup := j.isDuplicate(task.Username)
	if isDup {
		json.NewEncoder(w).Encode(Inside{Id: key})
		return
	} else {
		j.Create(user.Id, task.Username)
		j.JobsChan <- user.Id

		db.AddRecord(j.DB, user.Id, task.Username, statusInWork)
	}
	json.NewEncoder(w).Encode(user)
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

func (j *Worker) UpdateStatus(id string, ii ImportantInfo, state string) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id].State = state
	j.SafeZone.Status[id].Output = ii
}

func (j *Worker) Create(id string, username string) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id] = &OutputData{
		State:  statusInWork,
		Output: ImportantInfo{Username: username},
	}
}

func (j *Worker) SetErrStatus(id string, errCode int, description string) {
	j.SafeZone.Mu.Lock()
	defer j.SafeZone.Mu.Unlock()
	j.SafeZone.Status[id].State = StatusError
	j.SafeZone.Status[id].Error = &ErrInfo{Err: errCode, Description: description}
}

func (j *Worker) getUserInfo(id string) (data OutputData, notFound bool) {
	j.SafeZone.Mu.RLock()
	defer j.SafeZone.Mu.RUnlock()
	if val, ok := j.SafeZone.Status[id]; ok {
		data = *val
		return data, notFound
	} else {
		notFound = true
	}
	return data, notFound
}

func (j *Worker) GetUsernameById(id string) (username string) {
	j.SafeZone.Mu.RLock()
	defer j.SafeZone.Mu.RUnlock()
	if val, ok := j.SafeZone.Status[id]; ok {
		username = val.Output.Username
		return username
	}
	return username
}

// Check connection
func tryConnection(user string, pass string, database string) *gorm.DB {

	for errCounter := 0; errCounter < 6; errCounter++ {
		res, err := db.Init(user, pass, database)
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		return res
	}

	log.Error("Sorry but i cant connect to db try again")
	return nil
}

//If we want use db May be we will try create handler at here
func New(bufferSize int, user string, pass string, database string) *Worker {
	jobs := make(chan string, bufferSize)
	status := make(map[string]*OutputData)
	Db := tryConnection(user, pass, database)
	w := &Worker{JobsChan: jobs, SafeZone: SafeMapState{Status: status}, DB: Db}
	return w
}

func (j *Worker) Start(port string) {
	log.Info("Server listen ...")
	r := mux.NewRouter()
	r.HandleFunc("/users", j.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", j.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}/status", j.GetUserStatus).Methods("GET")
	r.HandleFunc("/users", j.ParseUser).Methods("POST")
	//r.HandleFunc("/users/{id}", w.PutUser).Methods("PUT")
	//r.HandleFunc("/users/{id}", w.DeleteUser).Methods("DELETE")

	// Не очень понял как именно закрытие сделать
	// Посмотрел пару примеров, но как то ...
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Some data but I don't know what is a data

}

func (j *Worker) Close() {
	time.Sleep(time.Second)
	db.CloseDb(j.DB)
}
