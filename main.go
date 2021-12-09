package main

import (
	"InstaSniffer/api"
	"InstaSniffer/info"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
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

		//check only index without User struct
		var index string
		for key, value := range w.SafeZone.Status {
			if value.State == "On Work" && value.Output.Username == w.SafeZone.Status[key].Output.Username {
				log.Info("Check key and value: ", key)
				index = key
			}
		}

		if len(j) == 0 {
			time.Sleep(time.Millisecond)
			log.Info("Worker ", id, " finished: ", j)
			w.Update(index, dummy)
			return
		}

		err, ii := info.UploadData(j)
		if err != nil {
			w.Update(index, dummy)
			log.Error(err)
			return
		}
		time.Sleep(time.Millisecond)
		log.Info("Worker ", id, " finished: ", j)
		w.Update(index, ii)
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

	// Create connection with api

	coreWorker()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)

	// TODO:
	// + 3. main.go 39, 58 - при вызове функций передаём туда w (api.Worker) - и при этом происходит копирование объекта мьютекса (mu).
	//		Это некорректно. Везде при вызове функций следует передавать указатель на api.Worker. Красиво было бы сделать так:
	// + 1. Поменять название метода PostUser на что-то более осмысленное (напр. ParseUser)
	// + 2. Если не может спарсить данные должен быть ответ dummy, но в ответе - пустая структура (надо проверить, почему так)
	// + 3. Сделать то что я писала по 3 пункту (либо пиши, что именно не получается)
	// + 4. Заменить все принты на логи (log.info..)
	// + 5. В пост запросе возвращать id (то есть генерить id на сервере, желательно uid).
	//		Соотвественно в теле запроса останется только name. Так это будет ближе к реальности.
	// + 6. Update, GetErrStatus - убрать из параметров state, т.к. мапа уже есть внутри j (Worker)
	// + 7. В Update сделать обновление, по аналогии с getErrStatus
	// + 8. Убрать глобальную переменную var users []User. Для выдачи статуса использовать мапу, где храним всю инфу о юзерах
	// + 9. Если ввожу несуществующий айди при проверке статуса - выдаётся статус 200 - всё ок (должна выводиться ошибка
	//		+ статус не 200).
	//		Плюс описание ошибки планировали выводить внутри OutputData, это уже где-то есть?
	// + 10* Не запускать задачу, если она уже в работе (Если уже есть в мапе не надо заново парсить)
	// * БД - сохранять результаты в таблицу users
	// 11 return Error struct in OutputData
}
