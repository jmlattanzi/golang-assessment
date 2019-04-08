package personcontroller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmlattanzi/golang-assessment/db"
	"github.com/jmlattanzi/golang-assessment/models"
)

// GetAllPeople ...Returns everyone
func GetAllPeople() []models.Person {
	people, _ := db.RunQuery("select * from address")
	return people
}

// GetSinglePerson ... Returns a single person with the id passed in
func GetSinglePerson(id string) models.Person {
	person, err := db.RunQuery("select * from address where id = " + id)
	if err != nil {
		log.Panicln("Error getting person: ", err)
	}

	return person[0]
}

// AddPerson ...Adds a new person to t he DB
func AddPerson(newPerson models.Person) models.Person {
	db.RunTransaction("insert into address (first_name, last_name, email, phone) values ($1, $2, $3, $4)", newPerson.FirstName, newPerson.LastName, newPerson.Email, newPerson.Email)
	return newPerson
}

// EditPerson ...Updates a person
func EditPerson(id string, updatedPerson models.Person) models.Person {
	db.RunTransaction("update address set first_name = $1, last_name = $2, email = $3, phone = $4 where id = "+id, updatedPerson.FirstName, updatedPerson.LastName, updatedPerson.Email, updatedPerson.Phone)

	person, _ := db.RunQuery("select * from address where id = " + id)
	return person[0]
}

// DeletePerson ...Delete's a person from the DB based on id
func DeletePerson(id string) {
	db.RunTransaction("delete from address where id = " + id)
	fmt.Println("[+] Person removed")
}

// ExportToCSV ... Exports the result of GetAllPeople to a CSV file
func ExportToCSV() {
	people := GetAllPeople()

	addressBook, err := os.Create("./address_book.csv")
	if err != nil {
		fmt.Println("Error creating file: ", err)
	}

	writer := csv.NewWriter(addressBook)
	for _, person := range people {
		var record []string
		record = append(record, strconv.Itoa(person.ID))
		record = append(record, person.FirstName)
		record = append(record, person.LastName)
		record = append(record, person.Email)
		record = append(record, person.Phone)
		writer.Write(record)
	}

	writer.Flush()
}

// ImportCSV ... Imports JSON address book to DB
func ImportCSV(people []byte) {
	person := []models.Person{}
	err := json.Unmarshal(people, &person)
	if err != nil {
		fmt.Println("[!] Error unmarshaling data: ", err)
	}

	for _, p := range person {
		db.RunTransaction("insert into address (first_name, last_name, email, phone) values ($1, $2, $3, $4)", p.FirstName, p.LastName, p.Email, p.Phone)
	}
}
