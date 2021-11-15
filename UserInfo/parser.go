package UserInfo

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadData(url string) error {

	fmt.Println(url)

	userData := Configure{}
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond * 10)

	req.Header.Set("User-Agent", "MSIE/15.0")

	res, err := spaceClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Cookies())

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &userData)
	if err != nil {
		return err
	}

	fmt.Printf("---------------------%T ----------------------", userData.Graphql)

	//Нужна структура полезной информации если так можно выразиться
	// Для каждого пользователя вот ее как раз и загоняем в файл

	//accounting("mua.shor")

	return nil
}

// Определится с полями какие в бд нужны и в каком виде

func accounting(userUid string) {
	file, err := os.Create(userUid + ".txt")
	if err != nil {
		log.Error(err)
	}
	defer file.Close()

	// json struct Info about user
	file.WriteString("Hello " + userUid)
}

func infoAboutUser() {

}
