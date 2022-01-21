package db

import (
	_ "database/sql"
	"time"
)

type dbData struct {
	name   string
	id     string
	status string
	date   time.Time
}

// Это выглядит плохо
//func Crete (DB *gorm.DB,info *api.OutputData, id string) error{
//	data := dbData{
//		name: info.Output.Name,
//		id: id,
//		status: info.State,
//		date: time.Now(),
//	}
//	DB.Create(data)
//}

func GetInfoById() {

}
