package main

import (
	"InstaSniffer/api"
	"InstaSniffer/info"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func worker(id int, w api.Worker, dummy api.ImportantInfo) {
	for j := range w.JobsChan {
		fmt.Println("Worker", id, "starting new:", j)
		err, ii := info.UploadData(j)
		if err != nil {
			w.Update(j, w.SafeZone.Status, dummy)
			log.Error(err)
		}
		time.Sleep(time.Millisecond)
		fmt.Println("Worker", id, "finished:", j)
		w.Update(j, w.SafeZone.Status, ii)
	}
}

func coreWorker() {
	// Буферизация на 10 элементов
	jobs := make(chan string, 10)

	//request status
	status := make(map[string]*api.OutputData)
	var mu sync.Mutex
	w := api.Worker{JobsChan: jobs, SafeZone: api.SafeMapState{Status: status, Mu: mu}}

	go api.ConnectionAPI(w)

	//dummy user
	// Info about model
	// Work only empty data in ImportantInfo
	dummy := api.ImportantInfo{
		Avatar:    "None",
		Name:      "TestValue",
		Username:  "test_value",
		Bio:       "None",
		CreatedAt: time.Now(),
	}

	//env
	thread := os.Getenv("THR")
	count, err := strconv.Atoi(thread)
	if err != nil {
		log.Error(err)
	}
	for i := 1; i <= count; i++ {
		go worker(i, w, dummy)
	}
}

//go:generate oapi-codegen -generate types -package api -o api/api.gen.go swagger.yaml

func main() {

	// Create connection with api

	coreWorker()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)

	// TODO:
	// +1. Реализовать роуты в апи (новая задача и получение результата)
	// 		1.1 Save res in map with sys_calls
	// +2. Настройки через переменные окружения (при запуске докера -e)
	// +3. Swagger (генерить структуры через go generate)
	// 4. Добавить коды ошибок в апи (404, 500...)
	// +5. Добавить дефолтного пользователя, если нет возможности авторизоваться
	// * БД - сохранять результаты в таблицу users
}
