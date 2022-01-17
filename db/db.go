package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	HOST = "localhost"
	PORT = 5432
)

type Database struct {
	Conn *sql.DB
}

//Input: User data (PostgreSQL script)
//Output: "Ok (all data about request)" or "Error: ..."
func Initialize(username, password, database string) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")
	return db, nil
}

type TestData struct {
	gorm.Model
	id         string
	name       string
	status     string
	created_at time.Time
}

func TestCreation() {
	sqlDB, err := sql.Open("postgres", "postgres")
	if err != nil {
		fmt.Errorf("Error: idk", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	gormDB.Create(&TestData{id: "100", name: "test", status: "done", created_at: time.Now()})
}
