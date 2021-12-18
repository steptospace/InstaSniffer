package main

import (
	"InstaSniffer/api"
	"InstaSniffer/info"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
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
	Threads    int `required:"true" envconfig:"THR"`
	BufferSize int `required:"true" envconfig:"BS"`
}

func startWork(id int, w *api.Worker) {
	for j := range w.JobsChan {
		log.Info("Worker: ", id, " starting new: ", j)
		if val, ok := w.GetUsernameById(j); !ok {
			time.Sleep(time.Millisecond)
			log.Info("Worker ", id, " finished: ", j)
			w.Update(j, dummy)
			return
		} else {
			err, ii := info.UploadData(val)
			if err != nil {
				w.Update(j, dummy)
				w.SetErrStatus(j, http.StatusNotFound, err.Error())
				log.Error(err)
				return
			}
			time.Sleep(time.Millisecond)
			log.Info("Worker ", id, " finished: ", j)
			w.Update(j, ii)
		}
	}
}

func coreWorker() {
	//env magic
	var env envConfig
	err := envconfig.Process("THR", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	//request status

	w := api.New(env.BufferSize)
	go w.Start()

	for i := 1; i <= env.Threads; i++ {
		go startWork(i, w)
	}

}

//go:generate oapi-codegen -generate types -package api -o api/api.gen.go swagger.yaml

func main() {
	//info.StartCommunicate("create table if not exist;")
	// Create connection with api

	coreWorker()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)

	// TODO:
	// * БД - сохранять результаты в таблицу users
	// 11 return Error struct in OutputData
}
