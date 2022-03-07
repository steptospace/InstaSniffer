package db

import (
	_ "database/sql"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

var dummy = Users{
	UId:      "None",
	Avatar:   "None",
	Name:     "None",
	Username: "None",
	Bio:      "None",
}


//Re:create this method
//We must push all data
func AddRecord(DB *gorm.DB, id string, uid string, state string) {
	var result []Tasks
	DB.Find(&result)

	data := Tasks{
		Id:        id,
		UserId:    uid,
		Status:    state,
		CreatedAt: time.Now(),
	}

	// We recreated method

	res := DB.Create(&data)
	if res.Error != nil {
		log.Error(res.Error)
		return
	}
}

func GetInfoById(DB *gorm.DB, id string) (result Users, err error) {
	// Go to Tasks table and get UserId
	var taskRes Tasks
	connectionCheck := DB.First(&taskRes, id)
	if connectionCheck.Error != nil {
		log.Error(connectionCheck.Error)
		return dummy, connectionCheck.Error
	}

	// Use UID --> check other info
	connectionCheck = DB.First(&result, taskRes.UserId)
	if connectionCheck.Error != nil {
		log.Error(connectionCheck.Error)
		return dummy, connectionCheck.Error
	}

	return result, nil
}

func DeleteRecord(DB *gorm.DB, id string) {

	// Не все так просто
	// Нужно еще удалить все связанные с ним элементы
	DB.Delete(Tasks{}, id)
}

func UpdateRecord(DB *gorm.DB, id string, status string) (result Tasks, err error) {
	res := DB.First(&result, id)
	if res.Error != nil {
		log.Error(res.Error)
		return result, res.Error
	}
	result.Status = status
	DB.Save(&result)
	return result, nil
}
