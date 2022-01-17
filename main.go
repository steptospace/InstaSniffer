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

	"github.com/golang-migrate/migrate/v4"
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
	Threads    int `required:"true" envconfig:"THR"`
	BufferSize int `required:"true" envconfig:"BS"`
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
	//Work with DataBase
	//db, err := sql.Open("postgres", "postgresql://postgres:admin@localhost:5432/postgres?sslmode=disable")
	//if err != nil {
	//	panic(err)
	//}
	//driver, err := postgres.WithInstance(db, &postgres.Config{})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(driver)
	m, err := migrate.New(
		"file://db/migration",
		"postgresql://postgres:admin@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Error(err)
	}
	// May be this one
	if err := m.Up(); err != migrate.ErrNoChange {
		log.Error(err)
	}

	// Create connection with api

	var env envConfig
	err = envconfig.Process("THR", &env)
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

	// TODO:
	// * БД - сохранять результаты в таблицу users
	// tests
}
