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

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	Threads    int `required:"true" envconfig:"THR" default:"2"`
	BufferSize int `required:"true" envconfig:"BS" default:"2"`
}

func startWork(id int, w *api.Worker) {
	for j := range w.JobsChan {
		log.Info("Worker: ", id, " starting new: ", j)
		name := w.GetUsernameById(j)
		err, ii := info.UploadData(name)
		if err != nil {
			w.SetErrStatus(j, http.StatusNotFound, err.Error())
			dummy.Username = name
			w.UpdateStatus(j, dummy, api.StatusError)
			log.Error(err)
			continue
		}
		time.Sleep(time.Millisecond)
		log.Info("Worker ", id, " finished: ", j)
		w.UpdateStatus(j, ii, api.StatusDone)
	}
}

//go:generate oapi-codegen -generate types -package api -o api/api.gen.go swagger.yaml

func main() {

	//Port int `required:"true" envconfig:"PORT" default:"8080"`
	//Пока не знаю как пробросить

	// Create connection with api
	var env envConfig
	err := envconfig.Process("THR", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	w := api.New(env.BufferSize)
	go w.Start()

	for i := 1; i <= env.Threads; i++ {
		go startWork(i, w)
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)
}
