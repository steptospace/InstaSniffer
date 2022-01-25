package api

import (
	"InstaSniffer/db"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateNewTask(name string) (id string) {

	testProfile := User{Username: name}
	data, err := json.Marshal(testProfile)
	if err != nil {
		fmt.Errorf("json.Marshal cant understand format")
	}
	res, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(data))
	if status := res.StatusCode; status != http.StatusOK {
		fmt.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Unknown data")
	}
	bodyString := string(bodyBytes)
	if len(bodyString) == 0 {
		fmt.Errorf("Empty data")
	}

	testId := Inside{}
	err = json.Unmarshal(bodyBytes, &testId)
	if err != nil {
		fmt.Errorf("Unmarshal cant use")
	}

	if testId.Id == "" {
		fmt.Errorf("Empty id columns")
	}
	return testId.Id
}

func TestAddNewUser(t *testing.T) {
	w := New(1, "postgres", "admin", "postgres")
	go w.Start("8080")
	// Create new user and check data
	id := CreateNewTask("rahmaninov")
	db.DeleteRecord(w.DB, id)
	//w.Close()
}

func TestGetEmptyInfo(t *testing.T) {
	w := New(1, "postgres", "admin", "postgres")
	go w.Start("8080")

	rr := httptest.NewRecorder()
	resp, err := http.Get("http://localhost:8080/users")
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unknown data")
	}
	bodyString := string(bodyBytes)

	var expected = "null\n"
	if bodyString != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUsers(t *testing.T) {
	w := New(1, "postgres", "admin", "postgres")
	go w.Start("8080")

	name := "rahmaninov"
	id := CreateNewTask(name)

	rr := httptest.NewRecorder()
	resp, err := http.Get("http://localhost:8080/users")
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unknown data")
	}
	bodyString := string(bodyBytes)
	// Check the response body is what we expect.
	data := []byte("[\"rahmaninov\"]\n")
	if bodyString != string(data) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(data))
	}
	db.DeleteRecord(w.DB, id)
}

func TestGetUserInfo(t *testing.T) {
	w := New(1, "postgres", "admin", "postgres")
	go w.Start("8080")
	name := "rahmaninov"
	id := CreateNewTask(name)

	resp, err := http.Get("http://localhost:8080/users/" + id)
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	testOutput := &OutputData{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unknown data")
	}
	err = json.Unmarshal(bodyBytes, &testOutput)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response body is what we expect.

	if testOutput.Output.Username != name {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Body, name)
	}
	db.DeleteRecord(w.DB, id)
}

func TestGetUserStatus(t *testing.T) {
	w := New(1, "postgres", "admin", "postgres")
	go w.Start("8080")
	name := "rahmaninov"

	id := CreateNewTask(name)

	resp, err := http.Get("http://localhost:8080/users/" + id + "/status")
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	testOutput := &OutputData{}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unknown data")
	}
	err = json.Unmarshal(bodyBytes, &testOutput)
	if err != nil {
		t.Fatal(err)
	}

	if testOutput.Output.Username != name {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Body, name)
	}
	db.DeleteRecord(w.DB, id)
}
