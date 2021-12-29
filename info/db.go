package info

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	user     = "postgres"
	password = "admin"
	host     = "localhost"
	port     = 5432
	dbname   = "postgres"
)

func Connect(userName string, password string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, userName, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error(err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Error(err)
		return nil
	}
	fmt.Println("User connection ...\n Success")
	return db
}

func Close(db *sql.DB) {
	db.Close()
}

//Input: User data (PostgreSQL script)
//Output: "Ok (all data about request)" or "Error: ..."
func StartCommunicate(db *sql.DB, textRequest string) (string, error) {
	start := time.Now()
	rows, err := db.Query(textRequest) // send command to database And execute
	if err != nil {
		log.Error(err)
	}

	/*if _, err := os.Stat(logPath + "\\logs.txt"); errors.Is(err, os.ErrNotExist) {
		log.Error(err)
	}
	*/
	end := time.Now()
	delta := end.Sub(start)

	fmt.Println(delta)
	// call any function
	info := ""
	//

	defer rows.Close()
	return info, nil
}
