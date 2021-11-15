package InstaSniffer

import (
	"InstaSniffer/API"
	"InstaSniffer/UserInfo"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func worker(id int, jobs <-chan string, result chan<- string) {
	for j := range jobs {
		fmt.Println("Worker", id, "starting new:", j)
		if err := UserInfo.UploadData(j); err != nil {
			log.Error(err)
		}
		time.Sleep(time.Millisecond)
		fmt.Println("Worker", id, "finished:", j)
		result <- j
		//передается структура надо запарсить только Нужные поля и сделать файлик
	}
}

func coreWorker() {

	// Work with this case
	jobList := []string{"steptospace", "mua.shor", "sweetheart_snail"}

	numJobs := len(jobList)
	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)

	//Позже поправить. Сделать через env
	for w := 1; w <= 1; w++ {
		go worker(w, jobs, results)
	}

	for _, login := range jobList {
		fmt.Println("Start job")
		jobs <- "http://www.instagram.com/" + login + "/?__a=1"
	}
	time.Sleep(time.Second * 10)
}

func main() {
	// Create connection with API
	API.ConnectionAPI()
	//coreWorker()
}
