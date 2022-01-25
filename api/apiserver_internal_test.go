package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddNewUser(t *testing.T) {
	w := New(1)
	go w.Start()
	// Create new user and check data
	testProfile := User{Username: "rahmaninov"}
	data, err := json.Marshal(testProfile)
	if err != nil {
		t.Error("json.Marshal cant understand format")
	}

	// request to api-server
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(w.ParseUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	testId := Inside{Id: "498081"}
	data, err = json.Marshal(testId)
	if err != nil {
		t.Fatal(err)
	}

	// On system programming 0x0a == \n
	data = append(data, 0x0a)

	if rr.Body.String() != string(data) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(data))
	}

	w.Close()
}

func TestGetEmptyInfo(t *testing.T) {
	w := New(1)
	go w.Start()

	_, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

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
	w := New(1)
	go w.Start()

	testProfile := User{Username: "rahmaninov"}
	data, err := json.Marshal(testProfile)
	if err != nil {
		t.Error("json.Marshal cant understand format")
	}

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(w.ParseUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	testId := Inside{Id: "498081"}
	data, err = json.Marshal(testId)
	if err != nil {
		t.Fatal(err)
	}

	// On system programming 0x0a == \n
	data = append(data, 0x0a)

	if rr.Body.String() != string(data) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(data))
	}

	req, err = http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
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
	data = []byte("[\"rahmaninov\"]\n")
	if bodyString != string(data) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(data))
	}
}

func TestGetUserInfo(t *testing.T) {
	w := New(1)
	go w.Start()
	data := `{
    "username": "rahmaninov"
}`
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(data)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(w.ParseUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":"498081"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	req, err = http.NewRequest("GET", "/users/{id}", nil)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"id": "498081",
	}

	rr = httptest.NewRecorder()
	req = mux.SetURLVars(req, vars)
	handler = http.HandlerFunc(w.GetUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected = "\"rahmaninov\"\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetUserStatus(t *testing.T) {
	w := New(1)
	go w.Start()

	testProfile := User{Username: "rahmaninov"}
	data, err := json.Marshal(testProfile)
	if err != nil {
		t.Error("json.Marshal cant understand format")
	}

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(w.ParseUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	testId := Inside{Id: "498081"}
	data, err = json.Marshal(testId)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, 0x0a)

	if rr.Body.String() != string(data) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(data))
	}

	req, err = http.NewRequest("GET", "/users/{id}/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	resp, err := http.Get("http://localhost:8080/users/498081/status")
	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	testOutput := &OutputData{}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Unknown data")
	}
	err = json.Unmarshal(bodyBytes, &testOutput)
	if err != nil {
		t.Fatal(err)
	}

	if testOutput.Output.Username != "rahmaninov" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "rahmaninov")
	}
}
