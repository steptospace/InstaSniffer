package main

import (
	"InstaSniffer/api"
	"InstaSniffer/info"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Work only when empty data in ImportantInfo
var dummy = api.ImportantInfo{
	Avatar:    "None",
	Name:      "None",
	Username:  "None",
	Bio:       "None",
	CreatedAt: time.Now(),
}

type envConfig struct {
	Threads int `required:"true" envconfig:"THR"`
}

func startWork(id int, w api.Worker) {
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

	//env magic
	var env envConfig
	err := envconfig.Process("THR", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	//request status
	status := make(map[string]*api.OutputData)
	var mu sync.Mutex
	w := api.Worker{JobsChan: jobs, SafeZone: api.SafeMapState{Status: status, Mu: mu}}

	go api.ConnectionAPI(w)

	for i := 1; i <= env.Threads; i++ {
		go startWork(i, w)
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
	// + 1. Заменить в свагере структуру респонза для /users/{id}/status
	// + 2. dummy можно вынести в var (снаружи функций), чтобы не передавать через параметры
	//3. main.go 39, 58 - при вызове функций передаём туда w (api.Worker) - и при этом происходит копирование объекта мьютекса (mu).
	//Это некорректно. Везде при вызове функций следует передавать указатель на api.Worker. Красиво было бы сделать так:
	// + 4. Получение переменных окружения желательно вынести в отдельную функцию, в идеале создать структуру (типа envConfig).
	//		И туда сохранить нужные переменные окружения (на случай, если их больше одной). В этой же функции весь парсинг переменных.
	//		Переменную count лучше назвать countWorkers (чтобы было сразу понятно о чём это)
	// + 5. Функцию worker лучше назвать startWork или как то глаголом (это по неймингу функций - должно быть действие)
	// + 6. parser.go - 56, 58 - не надо 2 раза вызывать infoAboutUser - лучше сохранить в переменную
	// + 7. parser.go - 72 - лучше сделать так:
	// + 8. Убрать ненужные роуты из core.go
	// + 9. Unlock() в мьютексах можно через defer делать
	// * БД - сохранять результаты в таблицу users
}
