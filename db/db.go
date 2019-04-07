package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmlattanzi/assessment/models"

	"github.com/jmlattanzi/assessment/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

// RunQuery ... Runs an SQL query
func RunQuery(q string) ([]models.Person, error) {
	config := config.LoadConfigurationFile("config.json")
	people := []models.Person{}

	// Open a new connection
	db, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		log.Println("[!] Error connecting to the db (db/Connect): ", err)
		return nil, err
	}

	fmt.Println("[+] Connected to the database")
	defer db.Close()

	rows, err := db.Query(q)
	if err != nil {
		log.Println("[!] Error querying db (db/RunQuery): ", err)
		return nil, err
	}

	for rows.Next() {
		person := models.Person{}
		err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Email, &person.Phone)
		if err != nil {
			fmt.Println("[!] Error scanning rows: ", err)
			return nil, err
		}

		people = append(people, person)
	}

	return people, nil
}

//RunTransaction ...Execute a SQL transaction
func RunTransaction(q string, args ...interface{}) error {
	fmt.Println("[-] Preparing SQL transaction")
	config := config.LoadConfigurationFile("config.json")
	db, err := sql.Open("postgres", config.DBURL)
	if err != nil {
		log.Println("[!] Error opening connection to db (RunTransaction): ", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("[!] Error beginning SQL transaction (RunTransaction): ", err)
		return err
	}

	stmt, err := tx.Prepare(q)
	if err != nil {
		log.Println("[!] Error preparing statement (RunTransaction): ", err)
		return err
	}

	fmt.Println("[+] SQL Transaction has been prepared")
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Println("[!] Error executing statement (RunTransaction): ", err)
		return err
	}

	tx.Commit()
	fmt.Println("[+] Transaction complete")
	return nil
}
