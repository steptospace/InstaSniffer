package db

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	dbURL := "postgres://postgres:admin@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	if err := db.AutoMigrate(&Data{}); err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDb(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(err)
		return
	}
	if err = sqlDB.Close(); err != nil {
		log.Error(err)
		return
	}
}
