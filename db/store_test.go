package db_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=postgres sslmode=disable"
	}
	os.Exit(m.Run())
}

func TestConnection (t *testing.T) {
}