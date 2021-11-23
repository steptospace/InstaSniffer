package main

import (
	"InstaSniffer/api"
	"InstaSniffer/info"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func worker(id int, jobs <-chan string, result chan<- string, status map[string]*api.OutputData) {
	for j := range jobs {
		fmt.Println("Worker", id, "starting new:", j)
		err, ii := info.UploadData(j)
		if err != nil {
			log.Error(err)
		}
		time.Sleep(time.Millisecond)
		fmt.Println("Worker", id, "finished:", j)

		// add mutex
		///sync.Mutex.Lock()
		tt := status[j]
		tt.State = "Done"
		tt.Output = ii
		///sync.Mutex.Unlock()
	}
}

func coreWorker() {
	// Буферизация на 10 элементов
	jobs := make(chan string, 10)
	results := make(chan string, 10)

	//request status
	status := make(map[string]*api.OutputData)

	w := api.Worker{JobsChan: jobs, Status: status}

	go api.ConnectionAPI(w)

	//env
	thread := os.Getenv("THR")
	count, err := strconv.Atoi(thread)
	if err != nil {
		log.Error(err)
	}
	for i := 1; i <= count; i++ {
		go worker(i, w.JobsChan, results, w.Status)
	}
}

func main() {

	// Create connection with api

	coreWorker()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)

	// TODO:
	// 1. Реализовать роуты в апи (новая задача и получение результата)
	// 1.1 Save res in map with sys_calls
	// +2. Настройки через переменные окружения (при запуске докера -e)
	// 3. Swagger (генерить структуры через go generate)
	// 4. Добавить коды ошибок в апи (404, 500...)
	// 5. Добавить дефолтного пользователя, если нет возможности авторизоваться
	// * БД - сохранять результаты в таблицу users
}
