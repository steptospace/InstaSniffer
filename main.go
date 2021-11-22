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

func worker(id int, jobs <-chan string, result chan<- string) {
	for j := range jobs {
		fmt.Println("Worker", id, "starting new:", j)
		if err := info.UploadData(j); err != nil {
			log.Error(err)
		}
		time.Sleep(time.Millisecond)
		fmt.Println("Worker", id, "finished:", j)
		result <- j
		//передается структура надо запарсить только Нужные поля и сделать файлик
	}
}

func coreWorker() {
	// Буферизация на 10 элементов
	jobs := make(chan string, 10)
	results := make(chan string, 10)

	fmt.Println("jobs", jobs, "res:", results)

	w := api.Worker{JobsChan: jobs}

	go api.ConnectionAPI(w)

	//env
	thread := os.Getenv("THR")
	count, err := strconv.Atoi(thread)
	if err != nil {
		log.Error(err)
	}
	for i := 1; i <= count; i++ {
		go worker(i, w.JobsChan, results)
	}
}

func main() {

	// Create connection with api

	coreWorker()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	for {
		fmt.Println("sleeping ...")
		time.Sleep(time.Second * 10)
	}

	// TODO:
	// 1. Реализовать роуты в апи (новая задача и получение результата)
	// 2. Настройки через переменные окружения (при запуске докера -e)
	// 3. Swagger (генерить структуры через go generate)
	// 4. Добавить коды ошибок в апи (404, 500...)
	// 5. Добавить дефолтного пользователя, если нет возможности авторизоваться
	// * БД - сохранять результаты в таблицу users
}
