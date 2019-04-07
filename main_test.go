package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jmlattanzi/assessment/controller"
	"github.com/jmlattanzi/assessment/models"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/person", controller.GetAllPeopleHandler).Methods("GET")
	router.HandleFunc("/person", controller.AddPersonHandler).Methods("POST")
	router.HandleFunc("/person/{id}", controller.GetSinglePersonHandler).Methods("GET")
	router.HandleFunc("/person/{id}", controller.EditPersonHandler).Methods("PUT")
	router.HandleFunc("/person/{id}", controller.DeletePersonHandler).Methods("DELETE")

	router.HandleFunc("/export", controller.ExportToCSVHandler).Methods("GET")
	router.HandleFunc("/import", controller.ImportCSVHandler).Methods("POST")

	return router
}

func TestGetAllPeopleHandler(t *testing.T) {
	fmt.Println("[t] Testing GetAllPeople...")
	req, _ := http.NewRequest("GET", "/person", nil)
	res := httptest.NewRecorder()

	Router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "OK response is expected")
}

func TestGetSinglePersonHandler(t *testing.T) {
	fmt.Println("[t] Testing GetSinglePersonHandler...")
	req, _ := http.NewRequest("GET", "/person/15", nil)
	res := httptest.NewRecorder()

	Router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "OK response expected")
}

func TestAddPersonHandler(t *testing.T) {
	fmt.Println("[t] Testing AddPersonHandler...")
	newPerson := &models.Person{
		FirstName: "Caleb",
		LastName:  "Horelius",
		Email:     "bigweeb@kawaii.co.jp",
		Phone:     "4537894789",
	}

	jsonPerson, _ := json.Marshal(newPerson)
	req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(jsonPerson))
	res := httptest.NewRecorder()

	Router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Should return a status of 200")
}

func TestEditPersonHandler(t *testing.T) {
	fmt.Println("[t] Testing EditPersonHandler")
	updatedPerson := &models.Person{
		FirstName: "Bruh",
		LastName:  "Moment",
	}

	jsonPerson, _ := json.Marshal(updatedPerson)
	req, err := http.NewRequest("PUT", "/person/15", bytes.NewBuffer(jsonPerson))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}
	res := httptest.NewRecorder()

	Router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Should be 200")
	// test if the updated person is returned
}

func TestDeletePersonHandler(t *testing.T) {
	fmt.Println("[t] Testing DeletePersonHandler")
	req, err := http.NewRequest("DELETE", "/person/10", nil)
	if err != nil {
		t.Error("Error creating a request: ", err)
	}
	res := httptest.NewRecorder()

	Router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Should be 200")
}

// should create a file
func TestExportToCSVHandler(t *testing.T) {
	fmt.Println("[t] Testing export to CSV")
	req, err := http.NewRequest("GET", "/export", nil)
	if err != nil {
		t.Error("Error creating a request", err)
	}

	res := httptest.NewRecorder()

	Router().ServeHTTP(res, req)

	if _, err := os.Stat("./address_book.csv"); err == nil {
		fmt.Println("File exists and was found")
	} else if os.IsNotExist(err) {
		t.Error("File not found: ", err)
	} else {
		t.Error("Something went wrong, see error: ", err)
	}

	assert.Equal(t, 200, res.Code, "Should be 200")
}

// THIS TEST WILL FAIL
// Not going to lie I've never setup a test for a request like this and I don't know how
// but now it's got me really curious about some concepts I haven't learned about
func TestImportCSVHandler(t *testing.T) {
	fmt.Println("[t] Testing ImportCSVHandler")

	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	defer pipeWriter.Close()
	_, err := writer.CreateFormFile("file", "./address_book.csv")
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", "/import", pipeReader)
	if err != nil {
		t.Error("[!] Error creating request: ", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Body = ioutil.NopCloser(pipeReader)

	res := httptest.NewRecorder()

	Router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Should be 200")
}
