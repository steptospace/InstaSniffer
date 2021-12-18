package info

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	user     = "postgres"
	password = "admin"
	host     = "localhost"
	port     = 5432
	dbname   = "postgres"
)

func createNewConnection() {
	fmt.Println("New connection")
}

func StartCommunicate(textRequest string) (string, error) {
	//Start conversation to db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return "Error: open connection dont work", err
		// Не забыть про RollBack!!!
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		return "Error: cant ping to DataBase", err
	}

	fmt.Println("Successfully connected!")

	// insert into * ()
	// Simply input
	rows, err := db.Query(textRequest) // send command to database
	if err != nil {
		//panic(err)
		return "Error: the request cannot be sent", err
	}

	defer rows.Close()

	return "Success", nil
}
