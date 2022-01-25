package db

import (
	_ "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func AddRecord(DB *gorm.DB, id string) {
	var result []Data
	DB.Find(&result)

	data := Data{
		Name:      "test name",
		Id:        id,
		CreatedAt: time.Now(),
		Status:    "done",
	}
	res := DB.Create(&data)
	if res.Error != nil {
		log.Error(res.Error)
		return
	}
}

func GetInfoById(DB *gorm.DB, id string) {
	var result Data
	res := DB.First(&result)
	if res.Error != nil {
		log.Error(res.Error)
		return
	}
	fmt.Println(res.Select(&result, id))
}
