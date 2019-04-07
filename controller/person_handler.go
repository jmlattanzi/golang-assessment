package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmlattanzi/assessment/models"
)

// GetAllPeopleHandler ...Returns all people from the DB
func GetAllPeopleHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	people := GetAllPeople()
	json.NewEncoder(res).Encode(&people)
}

// GetSinglePersonHandler ...Get a single person via ID
func GetSinglePersonHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	person := GetSinglePerson(params["id"])
	json.NewEncoder(res).Encode(&person)
}

// AddPersonHandler ... Adds a person to the DB
func AddPersonHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	newPerson := models.Person{}

	json.NewDecoder(req.Body).Decode(&newPerson)
	person := AddPerson(newPerson)
	json.NewEncoder(res).Encode(&person)
}

// EditPersonHandler ...Lets the user update a specific person
func EditPersonHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	updatedPerson := models.Person{}

	err := json.NewDecoder(req.Body).Decode(&updatedPerson)
	if err != nil {
		log.Fatal("[!] Error decoding body: ", err)
	}
	person := EditPerson(params["id"], updatedPerson)
	err = json.NewEncoder(res).Encode(&person)
	if err != nil {
		log.Fatal("[!] Error encoding person: ", err)
	}
}

// DeletePersonHandler ...Handles deleting a person from the DB
func DeletePersonHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	DeletePerson(params["id"])
	json.NewEncoder(res).Encode("You done it")
}

// ExportToCSVHandler ... Exports the entire DB to a .csv file
func ExportToCSVHandler(res http.ResponseWriter, req *http.Request) {
	ExportToCSV()
	http.ServeFile(res, req, "./address_book.csv")
}

// ImportCSVHandler ...Handles pulling file off the request body
func ImportCSVHandler(res http.ResponseWriter, req *http.Request) {
	person := models.Person{}
	people := []models.Person{}
	csvFile, _, err := req.FormFile("file")
	if err != nil {
		log.Fatal("[!] Error in FormFile: ", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1 // go through all of them

	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println("[!] Error reading csv data: ", err)
	}

	// add each person into the slice
	for _, p := range data {
		person.ID, _ = strconv.Atoi(p[0]) // convert the id to string
		person.FirstName = p[1]
		person.LastName = p[2]
		person.Email = p[3]
		person.Phone = p[4]
		people = append(people, person) // append them to the array
	}

	// marshal the data into json
	jsonPeople, err := json.Marshal(people)
	if err != nil {
		fmt.Println("[!] Error marshaling data: ", err)
	}

	// upload it to the db
	ImportCSV(jsonPeople)

	// respond with the encoded version of the slice
	json.NewEncoder(res).Encode(&people)
}
