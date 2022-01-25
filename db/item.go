package db

import (
	_ "database/sql"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func AddRecord(DB *gorm.DB, id string, name string, state string) {
	var result []Data
	DB.Find(&result)

	data := Data{
		Id:        id,
		Name:      name,
		Status:    state,
		CreatedAt: time.Now(),
	}
	res := DB.Create(&data)
	if res.Error != nil {
		log.Error(res.Error)
		return
	}
}

func GetInfoById(DB *gorm.DB, id string) (result Data, err error) {
	res := DB.First(&result, id)
	if res.Error != nil {
		log.Error(res.Error)
		return result, res.Error
	}
	return result, nil
}

func DeleteRecord(DB *gorm.DB, id string) {
	DB.Delete(Data{}, id)
}

func UpdateRecord(DB *gorm.DB, id string, status string) (result Data, err error) {
	res := DB.First(&result, id)
	if res.Error != nil {
		log.Error(res.Error)
		return result, res.Error
	}
	result.Status = status
	DB.Save(&result)
	return result, nil
}
